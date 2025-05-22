-- payback_migration.sql

-- Participants Table
CREATE TABLE IF NOT EXISTS participants (
    participant_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trips Table
CREATE TABLE IF NOT EXISTS trips (
    trip_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trip_Participants Junction Table
CREATE TABLE IF NOT EXISTS trip_participants (
    trip_id INTEGER NOT NULL,
    participant_id INTEGER NOT NULL,
    PRIMARY KEY (trip_id, participant_id),
    FOREIGN KEY (trip_id) REFERENCES trips(trip_id) ON DELETE CASCADE,
    FOREIGN KEY (participant_id) REFERENCES participants(participant_id) ON DELETE CASCADE
);

-- Original_Purchases Table
CREATE TABLE IF NOT EXISTS original_purchases (
    purchase_id INTEGER PRIMARY KEY AUTOINCREMENT,
    trip_id INTEGER NOT NULL,
    payer_participant_id INTEGER NOT NULL,
    total_amount INTEGER NOT NULL, -- In cents
    description TEXT NOT NULL,
    purchase_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (trip_id) REFERENCES trips(trip_id) ON DELETE CASCADE,
    FOREIGN KEY (payer_participant_id) REFERENCES participants(participant_id) ON DELETE RESTRICT -- Don't delete a participant if they've paid for something
);

-- Individual_Debts Table
CREATE TABLE IF NOT EXISTS individual_debts (
    debt_id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_purchase_id INTEGER NOT NULL,
    debtor_participant_id INTEGER NOT NULL,
    amount_owed INTEGER NOT NULL, -- In cents
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (original_purchase_id) REFERENCES original_purchases(purchase_id) ON DELETE CASCADE,
    FOREIGN KEY (debtor_participant_id) REFERENCES participants(participant_id) ON DELETE RESTRICT -- Don't delete a participant if they owe money
);

