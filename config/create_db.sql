CREATE TABLE IF NOT EXISTS commands 
(
    id_command         SERIAL PRIMARY KEY,
    script   TEXT NOT NULL,
    description_command TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS log_pids 
(
    id_pid SERIAL PRIMARY KEY,
    id_command  integer references commands(id_command),
    os_pid integer not null
);
CREATE TABLE IF NOT EXISTS data_pids 
(
    id_pid integer unique references log_pids(id_pid) ,
    data_start  timestamp,
    data_finish  timestamp
);

CREATE TABLE IF NOT EXISTS log_command 
(
    id SERIAL PRIMARY KEY,
    id_pid integer references log_pids(id_pid),
    data_logs   TEXT NOT NULL,
    type_log TEXT NOT NULL
);
