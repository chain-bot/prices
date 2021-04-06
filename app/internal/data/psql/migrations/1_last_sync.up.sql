CREATE TABLE last_sync(
    base_asset TEXT,
    quote_asset TEXT,
    exchange TEXT,
    last_sync TIMESTAMP WITH TIME ZONE,
    primary key (base_asset, quote_asset, exchange)
);