package accounts

import (
	"context"
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/bairn/infra/base"
	"resk/services"
	_ "resk/testx"
)

func NewAccountDomain() *accountDomain {
	return new(accountDomain)
}

type accountDomain struct {
	account Account
	accountLog AccountLog
}

func (domain *accountDomain) createAccountLogNo() {
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

func (domain *accountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

func (domain *accountDomain) createAccountLog() {
	domain.accountLog = AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo

	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.Username = domain.account.Username.String

	//交易对象信息
	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.Username.String

	domain.accountLog.Amount = domain.account.Balance
	domain.accountLog.Balance = domain.account.Balance

	//交易变化属性
	domain.accountLog.Decs = "账户创建"
	domain.accountLog.ChangeType = services.AccountCreated
	domain.accountLog.ChangeFlag = services.FlagAccountCreated

	domain.accountLog.Status = 1
}

func (domain *accountDomain) Create(dto services.AccountDTO) (*services.AccountDTO, error) {
	domain.account = Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.Username.Valid = true

	domain.createAccountLog()
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner

		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			return err
		}

		if id <= 0 {
			return errors.New("创建账户失败")
		}

		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}

		if id <= 0 {
			return errors.New("创建账户流水失败")
		}

		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	var rdto *services.AccountDTO
	rdto = domain.account.ToDTO()
	return rdto, err
}


func (a *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetOne(accountNo)
		return nil
	})

	if err != nil {
		return nil
	}

	if account == nil {
		return nil
	}

	return account.ToDTO()
}

//根据用户ID来查询红包账户信息
func (a *accountDomain) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(services.EnvelopeAccountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

//根据流水ID来查询账户流水
func (a *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetOne(logNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}

//根据交易编号来查询账户流水
func (a *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetByTradeNo(tradeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}

//根据用户ID和账户类型来查询账户信息
func (a *accountDomain) GetAccountByUserIdAndType(userId string, accountType services.AccountType) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(accountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

func (a *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status,err = a.TransferWithContextTx(ctx, dto)
		return err
	})
	return status, err
}

//必须在base.TX事务块里面运行，不能单独运行
func (a *accountDomain) TransferWithContextTx(ctx context.Context, dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	//如果交易变化是支出，修正amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	//创建账户流水记录
	a.accountLog = AccountLog{}
	a.accountLog.FromTransferDTO(&dto)
	a.createAccountLogNo()
	//检查余额是否足够和更新余额：通过乐观锁来验证，更新余额的同时来验证余额是否足够
	//更新成功后，写入流水记录
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}

		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = services.TransferedStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
			status = services.TransferedStatusSufficientFunds
			return errors.New("余额不足")
		}
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户出错")
		}
		a.account = *account
		a.accountLog.Balance = a.account.Balance
		id, err := accountLogDao.Insert(&a.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferedStatusFailure
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferedStatusSuccess
	}

	return status, err
}

