-- Initial schema for retail management system
-- Creates all base tables for the application

-- ============================================================================
-- CATEGORIA (Product Categories)
-- ============================================================================
CREATE TABLE IF NOT EXISTS categoria (
    id_categoria SERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL UNIQUE,
    descricao TEXT,
    data_criacao TIMESTAMP DEFAULT NOW(),
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================================
-- FORNECEDOR (Suppliers)
-- ============================================================================
CREATE TABLE IF NOT EXISTS fornecedor (
    id_fornecedor SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cnpj VARCHAR(14) NOT NULL UNIQUE,
    telefone VARCHAR(20),
    email VARCHAR(100),
    contato VARCHAR(100),
    
    -- Address
    endereco TEXT,
    cidade VARCHAR(50),
    estado VARCHAR(2),
    cep VARCHAR(8),
    
    -- Metadata
    data_cadastro TIMESTAMP DEFAULT NOW(),
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================================
-- USUARIO (System Users)
-- ============================================================================
CREATE TABLE IF NOT EXISTS usuario (
    id_usuario SERIAL PRIMARY KEY,
    nome VARCHAR(30) NOT NULL,
    nome_usuario VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    senha VARCHAR(256) NOT NULL,
    perfil VARCHAR(10) NOT NULL,
    ativo BOOLEAN NOT NULL,
    data_criacao TIMESTAMP DEFAULT NOW(),
    ultimo_login TIMESTAMP
);

-- ============================================================================
-- PRODUTO (Products)
-- ============================================================================
CREATE TABLE IF NOT EXISTS produto (
    id_produto SERIAL PRIMARY KEY,
    
    -- Product codes
    codigo_produto VARCHAR(30) UNIQUE NOT NULL, -- Internal SKU
    codigo_barras VARCHAR(14) UNIQUE,           -- EAN13 barcode
    
    -- Basic info
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    
    -- Relationships
    categoria_id INT NOT NULL,
    fornecedor_id INT,
    
    -- Pricing
    preco_custo NUMERIC(10,2) NOT NULL,
    preco_venda NUMERIC(10,2) NOT NULL,
    
    -- Inventory
    unidade_medida VARCHAR(10) DEFAULT 'UN',
    estoque_atual INT NOT NULL DEFAULT 0,
    estoque_minimo INT DEFAULT 0,
    controla_estoque BOOLEAN DEFAULT TRUE,
    
    -- Metadata
    data_criacao TIMESTAMP NOT NULL DEFAULT NOW(),
    data_atualizacao TIMESTAMP,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Foreign keys
    FOREIGN KEY (categoria_id) REFERENCES categoria(id_categoria),
    FOREIGN KEY (fornecedor_id) REFERENCES fornecedor(id_fornecedor)
);

-- ============================================================================
-- MOVIMENTACAO_ESTOQUE (Inventory Movements)
-- ============================================================================
CREATE TABLE IF NOT EXISTS movimentacao_estoque (
    id_movimentacao SERIAL PRIMARY KEY,
    produto_id INT NOT NULL,
    tipo_movimentacao VARCHAR(20) NOT NULL,
    quantidade INT NOT NULL,
    data_movimentacao TIMESTAMP NOT NULL DEFAULT NOW(),
    observacao TEXT,
    usuario_id INT,
    
    -- Foreign keys
    FOREIGN KEY (produto_id) REFERENCES produto(id_produto),
    FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario)
);

-- ============================================================================
-- CLIENTE (Customers)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cliente (
    id_cliente SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) UNIQUE,
    telefone VARCHAR(20),
    email VARCHAR(100),
    
    -- Address
    endereco TEXT,
    cidade VARCHAR(50),
    estado VARCHAR(2),
    cep VARCHAR(8),
    
    -- Metadata
    data_cadastro TIMESTAMP NOT NULL DEFAULT NOW(),
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================================
-- CAIXA (Cash Register)
-- ============================================================================
CREATE TABLE IF NOT EXISTS caixa (
    id_caixa SERIAL PRIMARY KEY,
    data_abertura TIMESTAMP NOT NULL,
    data_fechamento TIMESTAMP,
    valor_abertura NUMERIC(10,2) NOT NULL,
    valor_fechamento NUMERIC(10,2),
    status CHAR(1) NOT NULL,
    usuario_id INT NOT NULL,
    valor_sangria NUMERIC(10,2) DEFAULT 0,
    valor_suprimento NUMERIC(10,2) DEFAULT 0,
    
    -- Foreign keys
    FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario)
);

-- ============================================================================
-- VENDA (Sales)
-- ============================================================================
CREATE TABLE IF NOT EXISTS venda (
    id_venda SERIAL PRIMARY KEY,
    data_venda TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Pricing
    valor_bruto NUMERIC(10,2) NOT NULL,
    desconto NUMERIC(10,2) DEFAULT 0,
    valor_total NUMERIC(10,2) NOT NULL,
    
    -- Status and relationships
    status VARCHAR(20) NOT NULL,
    cliente_id INT,
    usuario_id INT NOT NULL,
    caixa_id INT NOT NULL,
    
    -- Foreign keys
    FOREIGN KEY (cliente_id) REFERENCES cliente(id_cliente),
    FOREIGN KEY (usuario_id) REFERENCES usuario(id_usuario),
    FOREIGN KEY (caixa_id) REFERENCES caixa(id_caixa)
);

-- ============================================================================
-- ITEM_VENDA (Sale Items)
-- ============================================================================
CREATE TABLE IF NOT EXISTS item_venda (
    id_item_venda SERIAL PRIMARY KEY,
    venda_id INT NOT NULL,
    produto_id INT NOT NULL,
    quantidade INT NOT NULL,
    preco_unitario NUMERIC(8,2) NOT NULL,
    subtotal NUMERIC(10,2) NOT NULL,
    custo_unitario NUMERIC(10,2) NOT NULL,
    
    -- Foreign keys
    FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
    FOREIGN KEY (produto_id) REFERENCES produto(id_produto)
);

-- ============================================================================
-- FORMA_PAGAMENTO (Payment Methods)
-- ============================================================================
CREATE TABLE IF NOT EXISTS forma_pagamento (
    id_forma_pagamento SERIAL PRIMARY KEY,
    descricao VARCHAR(15) NOT NULL,
    ativo BOOLEAN NOT NULL
);

-- ============================================================================
-- PAGAMENTO (Payments)
-- ============================================================================
CREATE TABLE IF NOT EXISTS pagamento (
    id_pagamento SERIAL PRIMARY KEY,
    venda_id INT NOT NULL,
    forma_pagamento_id INT NOT NULL,
    valor_pago NUMERIC(10,2) NOT NULL,
    data_pagamento TIMESTAMP DEFAULT NOW(),
    
    -- Foreign keys
    FOREIGN KEY (venda_id) REFERENCES venda(id_venda),
    FOREIGN KEY (forma_pagamento_id) REFERENCES forma_pagamento(id_forma_pagamento)
);
