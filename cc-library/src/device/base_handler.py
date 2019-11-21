class BaseHandler:

    def __init__(self):
        self.status = None


    def on_change(self, status):
        self.status = status