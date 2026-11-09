package model

// VeiculosPorLojaDto contagem de veículos disponíveis em cada concessionária.
type VeiculosPorLojaDto struct {
	LojaID     int    `json:"lojaId"`
	LojaNome   string `json:"lojaNome"`
	Quantidade int    `json:"quantidade"`
}

// PecaMaisVendidaDto item com alta saída nas vendas.
type PecaMaisVendidaDto struct {
	PecaID          int    `json:"pecaId"`
	PecaNome        string `json:"pecaNome"`
	TotalQuantidade int    `json:"totalQuantidade"`
}

// IndicadoresGerenciaisResponse consolidado gerencial retornado pela API.
type IndicadoresGerenciaisResponse struct {
	TotalVeiculosDisponiveis    int                  `json:"totalVeiculosDisponiveis"`
	TotalVeiculosReservados     int                  `json:"totalVeiculosReservados"`
	TotalVeiculosVendidos       int                  `json:"totalVeiculosVendidos"`
	TotalReservasAtivas         int                  `json:"totalReservasAtivas"`
	TotalAgendamentosPendentes  int                  `json:"totalAgendamentosPendentes"`
	TotalContatosNaoRespondidos int                  `json:"totalContatosNaoRespondidos"`
	TotalVendasVeiculos         int                  `json:"totalVendasVeiculos"`
	ValorTotalVendasVeiculos    float64              `json:"valorTotalVendasVeiculos"`
	VeiculosDisponiveisPorLoja  []VeiculosPorLojaDto `json:"veiculosDisponiveisPorLoja"`
	TopPecasMaisVendidas        []PecaMaisVendidaDto `json:"topPecasMaisVendidas"`
}
