
CREATE TYPE transaction_status AS ENUM (
  'initiated',
  'pending',
  'successful',
  'failed',
  'cancelled',
  'refunded'
);

