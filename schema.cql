-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS reduced_spatial WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a table
CREATE TABLE IF NOT EXISTS reduced_spatial.track_points (
  id text,
  t timestamp,
  X double,
  y double,
  z double,
  PRIMARY KEY(id, t)  
);
