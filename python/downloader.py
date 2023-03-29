import pandas as pd
import pandas.io.sql as pio
import psycopg2

# Офнул, чтобы pands не душил варнами
import warnings
warnings.filterwarnings('ignore')

class Downloader:

    def __init__(self, conn):
        self.conn = conn
        
    def get_pools_list_DB(self):
        skills = []
        query = '''
            SELECT * FROM "pools";
        '''
        self.skills_list = pio.read_sql_query(query, self.conn)
        return self.skills_list
    
    def get_dexes_list_DB(self):
        query = '''
            SELECT * FROM dexes;        
        '''
        return pio.read_sql_query(query, self.conn)

    def get_tokens_list_DB(self):
        query = '''
            SELECT * FROM tokens;        
        '''
        return pio.read_sql_query(query, self.conn)
    
    def get_pools_count(self):
        query = '''
            SELECT COUNT(*) FROM pools
        '''
        return pio.read_sql_query(query, self.conn)

    def get_dexes_count(self):
        query = '''
            SELECT COUNT(*) FROM dexes
        '''
        return pio.read_sql_query(query, self.conn)

    def get_tokens_count(self):
        query = '''
            SELECT COUNT(*) FROM tokens
        '''.format(network)
        return pio.read_sql_query(query, self.conn)

    
    def get_tokens_count_by_network(self, network):
        query = '''
            SELECT COUNT(*) network FROM tokens
            WHERE network='{}'
        '''.format(network)
        return pio.read_sql_query(query, self.conn).network[0]

    def get_reserves_by_block(self, network, blockNumber):
        query = '''
            SELECT * FROM reserves        
            WHERE block_number='{}';
        '''.format(blockNumber)
        return pio.read_sql_query(query, self.conn)
        
    