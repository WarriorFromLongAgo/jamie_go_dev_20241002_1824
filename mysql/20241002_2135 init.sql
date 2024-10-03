CREATE DATABASE IF NOT EXISTS workflow_management;
USE workflow_management;

CREATE TABLE workflow_info
(
    id            INT AUTO_INCREMENT PRIMARY KEY COMMENT 'workflow id',
    workflow_name VARCHAR(128)  NOT NULL,
    to_addr       varchar(64)   not null,
    description   varchar(1024) NOT NULL COMMENT 'workflow description',
    create_by     varchar(64)   not null comment 'create_by user_id',
    create_addr   varchar(64)   not null comment 'create_addr',
    created_time  datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    updated_by    varchar(64)   null comment 'updated_by user_id',
    updated_addr  varchar(64)   null comment 'updated_addr',
    updated_time  datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
) COMMENT 'workflow_info';

CREATE TABLE workflow_approve
(
    id           INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
    workflow_id  INT COMMENT 'workflow_info_id',
    approve_addr varchar(64) not null,
    status       ENUM ('approved', 'rejected') DEFAULT 'rejected' COMMENT 'approve status：approved/rejected, default rejected',
    approve_time datetime,
    create_by    varchar(64) not null comment 'create_by user_id',
    create_addr  varchar(64) not null comment 'create_addr',
    created_time TIMESTAMP                     DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    updated_by   varchar(64) null comment 'updated_by user_id',
    updated_addr varchar(64) null comment 'updated_addr',
    updated_time TIMESTAMP                     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
) COMMENT 'workflow_approve';

CREATE TABLE management
(
    id               INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
    name             VARCHAR(100) NOT NULL COMMENT 'name',
    permission_level ENUM ('none', 'partial', 'full') DEFAULT 'none' COMMENT 'permission_level: none/partial/full, default none',
    addr             varchar(64)  not null comment 'wallet addr',
    anvil_info       varchar(64)  NOT NULL COMMENT 'anvil info',
    create_by        varchar(64)  not null comment 'create_by user_id',
    create_addr      varchar(64)  not null comment 'create_addr',
    created_time     TIMESTAMP                        DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    updated_by       varchar(64)  null comment 'updated_by user_id',
    updated_addr     varchar(64)  null comment 'updated_addr',
    updated_time     TIMESTAMP                        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
) COMMENT 'management';

insert into management(name, permission_level, addr, anvil_info, create_by, create_addr)
    value ('anthn', 'full', '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', 'anvil 0', 0, '0x0');
insert into management(name, permission_level, addr, anvil_info, create_by, create_addr)
    value ('authz', 'full', '0x70997970C51812dc3A010C7d01b50e0d17dc79C8', 'anvil 1', 0, '0x0');
insert into management(name, permission_level, addr, anvil_info, create_by, create_addr)
    value ('test1', 'partial', '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC', 'anvil 2', 0, '0x0');
insert into management(name, permission_level, addr, anvil_info, create_by, create_addr)
    value ('test2', 'partial', '0x90F79bf6EB2c4f870365E785982E1f101E93b906', 'anvil 3', 0, '0x0');

CREATE TABLE workflow_configuration
(
    id          INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
    code        varchar(64)   not null,
    value       varchar(64)   not null,
    description varchar(1024) NOT NULL COMMENT 'workflow description'
) COMMENT 'workflow_configuration';
insert into workflow_configuration(id, code, value, description) value (null, 'finalize_num', 64, 'eth slot safe finalize');

# CREATE TABLE scheduled_log
# (
#     id             INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
#     execution_time TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP COMMENT 'execution_time',
#     status         ENUM ('success', 'failure') default 'failure' NOT NULL COMMENT 'success/false, default failure',
#     error_message  TEXT COMMENT 'error',
#     create_by      varchar(64)                                   not null comment 'create_by user_id',
#     create_addr    varchar(64)                                   not null comment 'create_addr',
#     created_time   TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
#     updated_by     varchar(64)                                   null comment 'updated_by user_id',
#     updated_addr   varchar(64)                                   null comment 'updated_addr',
#     updated_time   TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
# ) COMMENT 'scheduled_log';

CREATE TABLE token_info
(
    id               INT AUTO_INCREMENT PRIMARY KEY COMMENT 'token_info_id',
    token_name       VARCHAR(100) NOT NULL,
    token_symbol     VARCHAR(64)  NOT NULL,
    contract_address VARCHAR(64)  NOT NULL,
    decimals         int          not null default 18,
    create_by        varchar(64)  not null comment 'create_by user_id',
    create_addr      varchar(64)  not null comment 'create_addr',
    created_time     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    updated_by       varchar(64)  null comment 'updated_by user_id',
    updated_addr     varchar(64)  null comment 'updated_addr',
    updated_time     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
) COMMENT 'token_info';

insert into token_info(id, token_name, token_symbol, contract_address, decimals, create_by, create_addr) VALUE
    (null, 'Test_USDT', 'Test_USDT', '0x700b6A60ce7EaaEA56F065753d8dcB9653dbAD35', '6', '0', '0x0');

CREATE TABLE token_transfer_log
(
    id               INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
    token_info_id    INT                                   not null,
    from_address     VARCHAR(42)                           NOT NULL,
    to_address       VARCHAR(42)                           NOT NULL,
    amount           BIGINT UNSIGNED                       NOT NULL,
    status           ENUM ('pending', 'success', 'failed') not null DEFAULT 'pending',
    retry_count      INT                                   not null DEFAULT 0 COMMENT 'retry_count, default 0',
    transaction_hash VARCHAR(66)                           not null COMMENT 'tx hash',
    create_by        varchar(64)                           not null comment 'create_by user_id',
    create_addr      varchar(64)                           not null comment 'create_addr',
    created_time     TIMESTAMP                                      DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    updated_by       varchar(64)                           null comment 'updated_by user_id',
    updated_addr     varchar(64)                           null comment 'updated_addr',
    updated_time     TIMESTAMP                                      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_time'
)
    COMMENT 'token_transfer_log';

CREATE TABLE block_info
(
    id                bigint AUTO_INCREMENT PRIMARY KEY COMMENT 'block_id',
    block_hash        VARCHAR(128)    NOT NULL,
    block_parent_hash VARCHAR(128)    NOT NULL,
    block_number      BIGINT UNSIGNED NOT NULL,
    timestamp         TIMESTAMP       NOT NULL,
    rlp_bytes         VARCHAR(128)    NOT NULL,
    created_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    UNIQUE KEY (block_number) COMMENT 'block_number unique index',
    UNIQUE KEY (block_hash) COMMENT 'block_hash unique index'
) COMMENT 'block_info';

CREATE TABLE IF NOT EXISTS transaction_info
(
    id                bigint AUTO_INCREMENT PRIMARY KEY COMMENT 'tx_id',
    block_hash        VARCHAR(128)    NOT NULL,
    block_number      BIGINT UNSIGNED NOT NULL,
    tx_hash           VARCHAR(128)    NOT NULL,
    from_address      VARCHAR(64)     NOT NULL,
    to_address        VARCHAR(128)    NOT NULL,
    token_address     VARCHAR(128)    NOT NULL,
    gas_fee           bigint          NOT NULL,
    amount            bigint          NOT NULL comment 'tx transfer token/eth amount',
    status            SMALLINT        NOT NULL DEFAULT 0 comment 'receipt_statue,0=pending,1=success,2=failed',
    transaction_index bigint          NOT NULL comment 'tx index in one block',
    tx_type           SMALLINT        NOT NULL DEFAULT 0 comment 'tx type,0=Legacy,2=DynamicFee,3=Blob',
    created_time      TIMESTAMP                DEFAULT CURRENT_TIMESTAMP COMMENT 'created_time',
    UNIQUE KEY (tx_hash) COMMENT 'tx hash unique index',
    INDEX (block_number) COMMENT 'block_number index',
    INDEX (from_address) COMMENT 'from_address index',
    INDEX (to_address) COMMENT 'to_address index'
);