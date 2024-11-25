from datetime import date, time, datetime

from PyQt6.QtWidgets import QDialog
from ui.input_date import Ui_Dialog as Ui_Date
from ui.input_time import Ui_Dialog as Ui_Time
from ui.input_datetime import Ui_Dialog as Ui_DateTime


class DateDialog(QDialog, Ui_Date):

    def __init__(self, parent, title, default: date | None = None) -> None:
        super().__init__(parent)
        self.setWindowTitle(title)
        self.setupUi(self)
        if default is not None:
            self.dateEdit.setDate(default)


class TimeDialog(QDialog, Ui_Time):

    def __init__(self, parent, title, default: time | None = None) -> None:
        super().__init__(parent)
        self.setWindowTitle(title)
        self.setupUi(self)
        if default is not None:
            self.timeEdit.setTime(default)


class DateTimeDialog(QDialog, Ui_DateTime):

    def __init__(self,
                 parent,
                 title,
                 default: datetime | None = None,
                 label: str | None = None) -> None:
        super().__init__(parent)
        self.setWindowTitle(title)
        self.setupUi(self)
        if default is not None:
            self.dateTimeEdit.setDateTime(default)

        if label is not None:
            self.label.setText(label)
