PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

CREATE TABLE cluster(id integer primary key, name VARCHAR(100));
CREATE TABLE node(id integer primary key, name VARCHAR(100), cluster_id VARCHAR(200), network VARCHAR(15), mask int, FOREIGN KEY (cluster_id) REFERENCES cluster(id));
CREATE TABLE pod (id integer primary key, name VARCHAR(100), node_id int, FOREIGN KEY (node_id) REFERENCES node(id));
CREATE TABLE machine (id integer primary key, docker_id VARCHAR(300), pod_id int, FOREIGN KEY (pod_id) REFERENCES pod(id));
COMMIT;
