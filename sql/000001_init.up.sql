create table if not exists orders (
    order_uid          varchar(35) not null constraint orders_uid_pkey primary key,
    track_number       varchar(14),
    entry              varchar(5),
    locale             varchar(4),
    internal_signature varchar(255),
    customer_id        varchar(20),
    delivery_service   varchar(20),
    shardkey           varchar(20),
    sm_id              bigint,
    date_created       varchar(20),
    oof_shard          varchar(5)
);

create table if not exists items (
    order_uid    varchar(35) not null constraint items_orders_uid_fk references orders,
    chrt_id      bigint,
    track_number varchar(14),
    price        bigint,
    rid          varchar(21),
    name         varchar(255),
    sale         bigint,
    size         varchar(255),
    total_price  bigint,
    nm_id        bigint,
    brand        varchar(255),
    status       bigint
);

create table if not exists payments (
    order_uid     varchar(35) not null constraint payments_orders_id_fk references orders,
    transaction   varchar(35),
    request_id    varchar(255),
    currency      varchar(4),
    provider      varchar(15),
    amount        bigint,
    payment_dt    bigint,
    bank          varchar(255),
    delivery_cost bigint,
    goods_total   bigint,
    custom_fee    bigint
);

create table if not exists deliveries (
    order_uid varchar(35) not null constraint deliveries_orders_id_fk references orders,
    name     varchar(255),
    phone    varchar(16),
    zip      varchar(7),
    city     varchar(255),
    address  varchar(255),
    region   varchar(255),
    email    varchar(255)
);