CREATE TABLE IF NOT EXISTS Services (
    services_name VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (services_name)
);


CREATE TABLE IF NOT EXISTS UserSubscr (
    user_id VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    services_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,

    PRIMARY KEY (user_id, services_name, start_date),
    FOREIGN KEY (services_name) REFERENCES Services(services_name)
);
