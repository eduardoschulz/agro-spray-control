-- TODO atualizar alguns nomes e verificar se precisa UNIQUE e not null em alguns casos


CREATE TABLE usuarios (
    cpf VARCHAR(11) PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    nivel_permissao SMALLINT NOT NULL DEFAULT 1 CHECK (nivel_permissao BETWEEN 1 AND 3),
    password_hash TEXT NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
--    atualizado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE fazenda (
    cod SERIAL PRIMARY KEY,
    localizacao VARCHAR (255) NOT NULL,
    cpf_proprietario VARCHAR(11) NOT NULL,
    FOREIGN KEY (cpf_proprietario) REFERENCES usuarios(cpf)

);

CREATE TABLE area (
    cod SERIAL PRIMARY KEY, 
    tamanho INTEGER NOT NULL,
    fazenda_cod SERIAL,
    FOREIGN KEY (fazenda_cod) REFERENCES fazenda(cod)
);

CREATE TABLE produtos (
    cod SERIAL PRIMARY KEY,
    descricao VARCHAR(100) NOT NULL,
    fabricante VARCHAR(100) NOT NULL,
    compquimica VARCHAR(100) NOT NULL
);


CREATE TABLE lotes (
    cod SERIAL PRIMARY KEY,
    dtValidade DATE NOT NULL,
    cod_produto SERIAL,
    FOREIGN KEY (cod_produto) REFERENCES produtos(cod)
);


CREATE TABLE pulverizacao (
    cod SERIAL PRIMARY KEY,
    dtAplicacao DATE NOT NULL,
    cod_lote SERIAL,
    cpf_responsavel VARCHAR(11),
    FOREIGN KEY(cod_lote) REFERENCES lotes(cod),
    FOREIGN KEY(cpf_responsavel) REFERENCES usuarios(cpf)
);

CREATE TABLE pulverizacao_areas (
    cod_pulv INTEGER NOT NULL,
    cod_area INTEGER NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cod_pulv, cod_area),  -- Chave primária composta
    FOREIGN KEY (cod_pulv) REFERENCES pulverizacao(cod) ON DELETE CASCADE,
    FOREIGN KEY (cod_area) REFERENCES area(cod) ON DELETE CASCADE
);





