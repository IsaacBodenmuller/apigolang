CREATE TABLE IF NOT EXISTS produto (
	id_produto SERIAL PRIMARY KEY,
	nome VARCHAR(50) NOT NULL,
	descricao VARCHAR(100),
	preco DECIMAL(6,2) NOT NULL,
	estoque_atual INT NOT NULL,
	estoque_minimo INT,
	data_criacao TIMESTAMP NOT NULL,
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS movimentacao_estoque (
	id_movimentacao SERIAL PRIMARY KEY,
	tipo_movimentacao CHAR(1) NOT NULL,
	quantidade INT NOT NULL,
	data_movimentacao TIMESTAMP NOT NULL,
	observacao VARCHAR(50),
	produto_id INT NOT NULL,
	FOREIGN KEY (produto_id) REFERENCES produto(id_produto)
);

CREATE TABLE IF NOT EXISTS categoria (
	id_categoria SERIAL PRIMARY KEY,
	nome VARCHAR(20) NOT NULL,
	descricao VARCHAR(30) NOT NULL,
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS fornecedor (
	id_fornecedor SERIAL PRIMARY KEY,
	nome VARCHAR(50) NOT NULL,
	cnpj VARCHAR(14) NOT NULL,
	telefone VARCHAR(11),
	email VARCHAR(20),
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS cliente (
	id_cliente SERIAL PRIMARY KEY,
	nome VARCHAR(50) NOT NULL,
	cpf VARCHAR(11) NOT NULL,
	telefone VARCHAR(11),
	email VARCHAR(20),
	data_cadastro TIMESTAMP NOT NULL,
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS usuario (
	id_usuario SERIAL PRIMARY KEY,
	nome VARCHAR(30) NOT NULL,
	nome_usuario VARCHAR(20) NOT NULL,
	email VARCHAR(50) NOT NULL,
	senha VARCHAR(256) NOT NULL,
	perfil VARCHAR(10) NOT NULL,
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS caixa (
	id_caixa SERIAL PRIMARY KEY,
	data_abertura TIMESTAMP NOT NULL,
	data_fechamento TIMESTAMP NOT NULL,
	valor_abertura NUMERIC(10,2) NOT NULL,
	valor_fechamento NUMERIC(10,2) NOT NULL,
	status CHAR(1) NOT NULL,
	usuario_id INT NOT NULL,
	FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario)
);

CREATE TABLE IF NOT EXISTS venda (
	id_venda SERIAL PRIMARY KEY,
	data_venda TIMESTAMP NOT NULL,
	valor_total DECIMAL(8,2) NOT NULL,
	status CHAR(1) NOT NULL,
	cliente_id INT NOT NULL,
	usuario_id INT NOT NULL,
	caixa_id INT NOT NULL,
	FOREIGN KEY (cliente_id) REFERENCES cliente(id_cliente),
	FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario),
	FOREIGN KEY (caixa_id) REFERENCES caixa(id_caixa)
);

CREATE TABLE IF NOT EXISTS item_venda (
	id_item_venda SERIAL PRIMARY KEY,
	venda_id INT NOT NULL,
	produto_id INT NOT NULL,
	quantidade INT NOT NULL,
	preco_unitario NUMERIC(8,2) NOT NULL,
	subtotal NUMERIC(10,2) NOT NULL,
	FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
	FOREIGN KEY (produto_id) REFERENCES produto(id_produto)
);

CREATE TABLE IF NOT EXISTS forma_pagamento (
	id_forma_pagamento SERIAL PRIMARY KEY,
	descricao VARCHAR(15) NOT NULL,
	ativo BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pagamento (
	id_pagamento SERIAL PRIMARY KEY,
	venda_id INT NOT NULL,
	forma_pagamento_id INT NOT NULL,
	valor_pago NUMERIC (10,2) NOT NULL,
	FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
	FOREIGN KEY (forma_pagamento_id) REFERENCES forma_pagamento(id_forma_pagamento)
);
