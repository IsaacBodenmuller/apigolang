-- =============================
-- DROPAR TUDO (caso já exista)
-- =============================
DROP TABLE IF EXISTS pagamento;
DROP TABLE IF EXISTS item_venda;
DROP TABLE IF EXISTS venda;
DROP TABLE IF EXISTS caixa;
DROP TABLE IF EXISTS movimentacao_estoque;
DROP TABLE IF EXISTS produto;
DROP TABLE IF EXISTS forma_pagamento;
DROP TABLE IF EXISTS cliente;
DROP TABLE IF EXISTS fornecedor;
DROP TABLE IF EXISTS categoria;
DROP TABLE IF EXISTS usuario;


-- =============================
-- TABELAS BASE
-- =============================

CREATE TABLE categoria (
	id_categoria SERIAL PRIMARY KEY,
	nome VARCHAR(50) NOT NULL UNIQUE,
	descricao TEXT,
	data_criacao TIMESTAMP DEFAULT NOW(),
	ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE fornecedor (
	id_fornecedor SERIAL PRIMARY KEY,
	nome VARCHAR(100) NOT NULL,
	cnpj VARCHAR(14) NOT NULL UNIQUE,
	telefone VARCHAR(20),
	email VARCHAR(100),
	contato VARCHAR(100),
	endereco TEXT,
	cidade VARCHAR(50),
	estado VARCHAR(2),
	cep VARCHAR(8),
	data_cadastro TIMESTAMP DEFAULT NOW(),
	ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE cliente (
	id_cliente SERIAL PRIMARY KEY,
	nome VARCHAR(100) NOT NULL,
	cpf VARCHAR(11) UNIQUE,
	telefone VARCHAR(20),
	email VARCHAR(100),
	endereco TEXT,
	cidade VARCHAR(50),
	estado VARCHAR(2),
	cep VARCHAR(8),
	data_cadastro TIMESTAMP NOT NULL DEFAULT NOW(),
	ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE usuario (
	id_usuario SERIAL PRIMARY KEY,
	nome VARCHAR(30) NOT NULL,
	nome_usuario VARCHAR(20) NOT NULL,
	email VARCHAR(50) NOT NULL,
	senha VARCHAR(256) NOT NULL,
	perfil VARCHAR(20) NOT NULL,
	role VARCHAR(10) NOT NULL,
	ativo BOOLEAN NOT NULL,
);

-- =============================
-- PRODUTO
-- =============================

CREATE TABLE produto (
	id_produto SERIAL PRIMARY KEY,
	codigo_produto VARCHAR(30) UNIQUE NOT NULL,
	codigo_barras VARCHAR(14) UNIQUE,
	nome VARCHAR(100) NOT NULL,
	descricao TEXT,
	categoria_id INT NOT NULL,
	fornecedor_id INT,
	preco_custo NUMERIC(10,2) NOT NULL,
	preco_venda NUMERIC(10,2) NOT NULL,
	unidade_medida VARCHAR(10) DEFAULT 'UN',
	estoque_atual INT NOT NULL DEFAULT 0,
	estoque_minimo INT DEFAULT 0,
	controla_estoque BOOLEAN DEFAULT TRUE,
	data_criacao TIMESTAMP NOT NULL DEFAULT NOW(),
	data_atualizacao TIMESTAMP,
	ativo BOOLEAN NOT NULL DEFAULT TRUE,
	FOREIGN KEY (categoria_id) REFERENCES categoria(id_categoria),
	FOREIGN KEY (fornecedor_id) REFERENCES fornecedor(id_fornecedor)
);

-- =============================
-- CAIXA
-- =============================

CREATE TABLE caixa (
	id_caixa SERIAL PRIMARY KEY,
	data_abertura TIMESTAMP NOT NULL,
	data_fechamento TIMESTAMP,
	valor_abertura NUMERIC(10,2) NOT NULL,
	valor_fechamento NUMERIC(10,2),
	status CHAR(1) NOT NULL,
	usuario_id INT NOT NULL,
	valor_sangria NUMERIC(10,2) DEFAULT 0,
	valor_suprimento NUMERIC(10,2) DEFAULT 0,
	FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario)
);

-- =============================
-- VENDA
-- =============================

CREATE TABLE venda (
	id_venda SERIAL PRIMARY KEY,
	data_venda TIMESTAMP NOT NULL DEFAULT NOW(),
	valor_bruto NUMERIC(10,2) NOT NULL,
	desconto NUMERIC(10,2) DEFAULT 0,
	valor_total NUMERIC(10,2) NOT NULL,
	status VARCHAR(20) NOT NULL,
	cliente_id INT,
	usuario_id INT NOT NULL,
	caixa_id INT NOT NULL,
	FOREIGN KEY (cliente_id) REFERENCES cliente(id_cliente),
	FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario),
	FOREIGN KEY (caixa_id) REFERENCES caixa(id_caixa)
);

-- =============================
-- ITEM VENDA
-- =============================

CREATE TABLE item_venda (
	id_item_venda SERIAL PRIMARY KEY,
	venda_id INT NOT NULL,
	produto_id INT NOT NULL,
	quantidade INT NOT NULL,
	preco_unitario NUMERIC(10,2) NOT NULL,
	subtotal NUMERIC(10,2) NOT NULL,
	custo_unitario NUMERIC(10,2) NOT NULL,
	FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
	FOREIGN KEY (produto_id) REFERENCES produto(id_produto)
);

-- =============================
-- PAGAMENTOS
-- =============================

CREATE TABLE forma_pagamento (
	id_forma_pagamento SERIAL PRIMARY KEY,
	descricao VARCHAR(15) NOT NULL,
	ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pagamento (
	id_pagamento SERIAL PRIMARY KEY,
	venda_id INT NOT NULL,
	forma_pagamento_id INT NOT NULL,
	valor_pago NUMERIC (10,2) NOT NULL,
	data_pagamento TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
	FOREIGN KEY (forma_pagamento_id) REFERENCES forma_pagamento(id_forma_pagamento)
);

-- =============================
-- MOVIMENTAÇÃO DE ESTOQUE
-- =============================

CREATE TABLE movimentacao_estoque (
	id_movimentacao SERIAL PRIMARY KEY,
	produto_id INT NOT NULL,
	tipo_movimentacao VARCHAR(20) NOT NULL,
	quantidade INT NOT NULL,
	data_movimentacao TIMESTAMP NOT NULL DEFAULT NOW(),
	observacao TEXT,
	usuario_id INT,
	FOREIGN KEY (produto_id) REFERENCES produto(id_produto),
	FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario)
);
