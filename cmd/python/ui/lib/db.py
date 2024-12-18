import psycopg2


class Database:

    def __init__(
        self, db_name: str, db_pass: str, db_user: str, db_host: str, db_port: int
    ) -> None:
        self.conn = psycopg2.connect(
            f"host={db_host} user={db_user} dbname={db_name} port={db_port} sslmode=disable password={db_pass}",
        )

    def get_info(self, table: str):
        with self.conn.cursor() as cur:
            cur.execute(f"SELECT * FROM {table}")  # trunk-ignore(bandit)
            headers = [desc[0] for desc in cur.description]

            return (headers, cur.fetchall())

    def update_data(self, table: str, value, id, col_name: str):
        with self.conn.cursor() as cur:
            cur.execute(
                f"UPDATE {table} SET {col_name} = % WHERE id = %s", (value, id) # trunk-ignore(bandit)
            )
            self.conn.commit()

    def __del__(self):
        self.conn.close()
