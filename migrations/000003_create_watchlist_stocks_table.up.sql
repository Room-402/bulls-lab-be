CREATE TABLE IF NOT EXISTS watchlist_stocks (
    id            SERIAL PRIMARY KEY,
    active        BOOLEAN NOT NULL DEFAULT true,
    stock_ticker  VARCHAR NOT NULL,
    watchlist_id  INTEGER NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE
);

CREATE INDEX        IF NOT EXISTS idx_watchlist_stocks_watchlist_id ON watchlist_stocks(watchlist_id);