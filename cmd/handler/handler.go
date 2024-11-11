package handler

import (
	transactionProto "github.com/justIGreK/MoneyKeeper-Transaction/pkg/go/transaction"
	"google.golang.org/grpc"
)

type Handler struct {
	server      grpc.ServiceRegistrar
	transaction TransactionService
}

func NewHandler(grpcServer grpc.ServiceRegistrar, txSRV TransactionService) *Handler {
	return &Handler{server: grpcServer, transaction: txSRV}
}
func (h *Handler) RegisterServices() {
	h.registerTxService(h.server, h.transaction)
}

func (h *Handler) registerTxService(server grpc.ServiceRegistrar, tx TransactionService) {
	transactionProto.RegisterTransactionServiceServer(server, &TransactionServiceServer{TxSRV: tx})
}
