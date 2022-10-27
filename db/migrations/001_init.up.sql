CREATE TABLE IF NOT EXISTS transactions (
  transaction_id INTEGER PRIMARY KEY,
  request_id INTEGER NOT NULL,
  terminal_id INTEGER NOT NULL,
  partner_object_id INTEGER NOT NULL,
  amount_total INTEGER NOT NULL, -- insert *100 get /100
  amount_original INTEGER NOT NULL, -- insert *100 get /100
  commission_ps INTEGER NOT NULL, -- insert *100 get /100
  commission_client INTEGER NOT NULL, -- insert *100 get /100
  commission_provider INTEGER NOT NULL, -- insert *100 get /100
  date_input TIMESTAMP NOT NULL,
  date_post TIMESTAMP NOT NULL,
  status TEXT NOT NULL,
  payment_type TEXT NOT NULL,
  payment_number TEXT NOT NULL,
  service_id INTEGER NOT NULL,
  service TEXT NOT NULL,
  payee_id INTEGER NOT NULL,
  payee_name TEXT NOT NULL,
  payee_bank_mfo INTEGER NOT NULL,
  payee_bank_account TEXT NOT NULL,
  payment_narrative TEXT NOT NULL
);
