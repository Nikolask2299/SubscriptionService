CREATE TABLE IF NOT EXISTS Services (
    services_name VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (services_name)
);

CREATE TABLE IF NOT EXISTS UserSubscr (
    subscrb_id SERIAL NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,   
    end_date DATE NULL,
    services_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,

    PRIMARY KEY (user_id, services_name, start_date),
    FOREIGN KEY (services_name) REFERENCES Services(services_name)
);

CREATE INDEX IF NOT EXISTS idx_subscr_id ON UserSubscr (subscrb_id);
CREATE INDEX IF NOT EXISTS idx_user_id_services ON UserSubscr (user_id, services_name);
CREATE INDEX IF NOT EXISTS idx_services_name_date ON UserSubscr (services_name, start_date);
CREATE INDEX IF NOT EXISTS idx_user_id_date ON UserSubscr (user_id, start_date);
