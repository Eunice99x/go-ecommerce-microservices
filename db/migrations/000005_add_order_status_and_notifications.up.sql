ALTER TABLE orders
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'shipped', 'delivered'));

-- outbox for customer emails: rows are written in the same transaction as the
-- order change and drained by the notifier. locked_until is a lease so two
-- workers never send the same row, and a crashed send is retried once it expires.
CREATE TABLE notifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    order_status VARCHAR(20) NOT NULL,
    state VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (state IN ('pending', 'sent', 'failed')),
    attempts INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ
);

CREATE INDEX notifications_pending_idx ON notifications (created_at) WHERE state = 'pending';
