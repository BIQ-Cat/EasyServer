import os
import json
import typing
import datetime
import decimal
import uuid

from PyQt6.QtCore import Qt
from PyQt6.QtWidgets import QWidget, QTreeWidgetItem, QTableWidgetItem, QInputDialog, QMessageBox
from psycopg2 import Error
from ui.lib.db import Database
from ui.lib import date_dialogs
from ui.frame import Ui_Form
from setup import GetDefaultModuleConfiguration, GetEnvironmentConfiguration

STR = 8
INT = 4
FLOAT = 0
BOOL = 2
DICT = 1
LIST = 16
DATE = 32

VALUE_TYPES = {STR: 0, BOOL: 1, INT: 2, DICT: -1, LIST: -4, FLOAT: 4, DATE: 8}


# Object should have fields `value` and `valueType`
def set_value_type(obj):
    if isinstance(obj.value, str):
        obj.valueType = VALUE_TYPES[STR]
    elif isinstance(obj.value, int):
        obj.valueType = VALUE_TYPES[INT]
    elif isinstance(obj.value, float):
        obj.valueType = VALUE_TYPES[FLOAT]
    elif isinstance(obj.value, dict):
        obj.valueType = VALUE_TYPES[DICT]
    elif isinstance(obj.value, list):
        obj.valueType = VALUE_TYPES[LIST]
    elif isinstance(obj.value, datetime.datetime):
        obj.valueType = VALUE_TYPES[DATE]

    return obj


class JsonField(QTreeWidgetItem):

    def __init__(self, parent: QWidget, file: str, path: list[str | int],
                 value: typing.Any, partIndex: int, aspectIndex: int):
        self.file = file
        self.path = path
        self.value = value
        self.valueType = VALUE_TYPES[STR]
        self.partIndex = partIndex
        self.aspectIndex = aspectIndex
        self.currentTable = None
        super().__init__(parent)
        self.setText(0, path[-1])
        print(GetEnvironmentConfiguration())

    def createChild(self, field: str, value: typing.Any) -> typing.Self:
        return JsonField(self, self.file, [*self.path, field], value,
                         self.partIndex, self.aspectIndex)

    def __str__(self) -> str:
        field_name = ''
        for el in reversed(self.path):
            if isinstance(el, int):
                field_name = '[' + str(el) + ']' + field_name
            else:
                field_name = el + field_name
                break
        else:
            field_name = os.path.splitext(os.path.basename(
                self.file))[0] + field_name

        return field_name


class DBField(QTableWidgetItem):

    def __init__(self, value: typing.Any):
        self.value = value
        super().__init__(str(self.value))
        self.setFlags(Qt.ItemFlag.ItemIsEnabled | Qt.ItemFlag.ItemIsSelectable)

    def __str__(self) -> str:
        return str(self.value)


class Project(QWidget, Ui_Form):

    def __init__(self, dir):
        self.dir = dir
        self.configPath = os.path.join(dir, 'EasyConfig')
        self.db = Database("easy_server", "root", "admin", "localhost", 5432)
        self.currentTable = ''
        super().__init__()
        self.setupUi(self)

        self.setupConfiguration()
        self.setupDatabase()

    def setupConfiguration(self):
        if not os.path.exists(self.configPath):
            os.mkdir(self.configPath)

        for i in range(self.treeWidget.topLevelItemCount()):
            part = self.treeWidget.topLevelItem(i)
            if not (os.path.exists(
                    os.path.join(self.configPath,
                                 part.text(0).lower()))):
                os.mkdir(os.path.join(self.configPath, part.text(0).lower()))

            for j in range(part.childCount()):
                aspect = part.child(j)
                filepath = os.path.join(self.configPath,
                                        part.text(0).lower(),
                                        aspect.text(0).lower() + '.json')
                if not os.path.exists(filepath):
                    with open(filepath, 'wb') as f:
                        content, ok = GetDefaultModuleConfiguration(
                            aspect.text(0).lower())
                        if not ok:
                            content = json.dumps(self.buildDict(aspect),
                                                 indent=2).encode()

                        f.write(content)

                with open(filepath) as f:
                    data = json.load(f)
                    self.updateAspect(i, j, data, filepath)

        self.treeWidget.itemDoubleClicked.connect(self.changeField)

    def updateAspect(self, partIndex: int, aspectIndex: int, data: dict,
                     file: str):
        part = self.treeWidget.topLevelItem(partIndex)
        aspect = part.child(aspectIndex)

        new_aspect = QTreeWidgetItem(part)
        new_aspect.setText(0, aspect.text(0))
        new_aspect = self.buildTree(data, new_aspect, file, [], partIndex,
                                    aspectIndex)
        part.removeChild(aspect)
        part.insertChild(aspectIndex, new_aspect)

    def changeField(self, item: QTreeWidgetItem, column):
        if not isinstance(item, JsonField):
            return

        ok = False
        value = None
        if item.valueType == VALUE_TYPES[STR]:
            value, ok = QInputDialog.getText(self,
                                             str(item),
                                             'Изменить текст:',
                                             text=item.value)
        elif item.valueType == VALUE_TYPES[INT]:
            value, ok = QInputDialog.getInt(self, str(item), 'Изменить число:',
                                            item.value)
        elif item.valueType == VALUE_TYPES[FLOAT]:
            value, ok = QInputDialog.getDouble(self, str(item),
                                               'Изменить число:', item.value)
        elif item.valueType == VALUE_TYPES[BOOL]:
            value, ok = QInputDialog.getItem(self, str(item),
                                             'Изменить значение:',
                                             ['false', 'true'],
                                             int(item.value))
            if ok and value:
                value = value == 'true'

        if ok and value:
            part = self.treeWidget.topLevelItem(item.partIndex)
            aspect = part.child(item.aspectIndex)
            data = self.buildDict(aspect)
            current_el = data
            for p in item.path[:-1]:
                if (current_el := current_el.get(p)) is None:
                    break
            else:
                current_el[item.path[-1]] = value

            with open(item.file, 'w') as f:
                json.dump(data, f, indent=2)

            self.updateAspect(item.partIndex, item.aspectIndex, data,
                              item.file)

    def buildDict(self, root: QTreeWidgetItem):
        data = {}
        for i in range(root.childCount()):
            el = root.child(i)
            if not el.text(1):
                data[el.text(0)] = self.buildDict(el)
                continue

            try:
                data[el.text(0)] = int(el.text(1))
            except ValueError:
                if el.text(1) == 'true' or el.text(1) == 'false':
                    data[el.text(0)] = el.text(1) == 'true'
                else:
                    data[el.text(0)] = el.text(1)

        return data

    def setValueToField(self, el: JsonField) -> JsonField:
        if isinstance(el.value, bool):
            el.valueType = VALUE_TYPES[BOOL]
            el.setText(1, 'true' if el.value else 'false')
        elif isinstance(el.value, dict):
            el.valueType = VALUE_TYPES[DICT]
            el = self.buildTree(el.value, el, el.file, el.path, el.partIndex,
                                el.aspectIndex)
        elif isinstance(el.value, int):
            el.valueType = VALUE_TYPES[INT]
            el.setText(1, str(el.value))
        elif isinstance(el.value, float):
            el.valueType = VALUE_TYPES[FLOAT]
            el.setText(1, str(el.value))
        elif isinstance(el.value, list):
            el.valueType = VALUE_TYPES[LIST]
            el = self.buildListInTree(el.value, el.root)
        else:
            el.setText(1, str(el.value))

        return el

    def buildListInTree(self, data: list, root: JsonField):
        for i, value in enumerate(data):
            el = root.createChild(str(i), value)
            el.path[-1] = i
            el = self.setValueToField()
            root.addChild(el)

    def buildTree(self, data: dict, root: QTreeWidgetItem, file: str,
                  path: list[str], partIndex: int, aspectIndex: int):
        for i in range(root.childCount()):
            root.removeChild(root.child(i))

        for key, value in data.items():
            el = self.setValueToField(
                JsonField(root, file, [*path, key], value, partIndex,
                          aspectIndex))
            root.addChild(el)

        return root

    def setupDatabase(self):
        self.tableChooser.currentIndexChanged.connect(self.changeTable)
        self.tableWidget.cellDoubleClicked.connect(self.editCell)

    def changeTable(self, i):
        self.currentTable = self.tableChooser.itemText(
            i).lower() if i != 0 else None
        self.updateTable()

    def updateTable(self):
        if self.currentTable is None:
            self.tableWidget.setRowCount(0)
            self.tableWidget.setColumnCount(0)
            return
        headers, data = self.db.get_info(self.currentTable)
        headers = list(map(str.capitalize, headers))
        try:
            headers[headers.index('Id')] = 'ID'
        except ValueError:
            pass

        self.tableWidget.setColumnCount(len(headers))
        self.tableWidget.setRowCount(len(data))

        self.tableWidget.setHorizontalHeaderLabels(headers)
        for row, line in enumerate(data):
            for col, value in enumerate(line):
                el = DBField(value)
                self.tableWidget.setItem(row, col, el)

    def editCell(self, row: int, col: int):
        el = self.tableWidget.item(row, col)
        if not isinstance(el, DBField) or col == 0:
            return

        title = 'Настройка базы данных'

        ok = False
        value = None
        if isinstance(el.value, bool):
            value, ok = QInputDialog.getItem(self, title, 'Изменить значение:',
                                             ['false', 'true'], int(el.value))
            if ok and value:
                value = value == 'true'
        elif isinstance(el.value, float):
            value, ok = QInputDialog.getDouble(self, title, 'Изменить число:',
                                               el.value)
        elif isinstance(el.value, int):
            value, ok = QInputDialog.getInt(self, title, 'Изменить число:',
                                            el.value)
        elif isinstance(el.value, decimal.Decimal):
            value, ok = QInputDialog.getText(
                self,
                title,
                "Изменить число (ВНИМАНИЕ! ВВОДИТЬ МОЖНО ТОЛЬКО ЧИСЛА!)",
                text=str(el.value))
            if ok and value:
                try:
                    new_value = decimal.Decimal(value)
                except decimal.InvalidOperation:
                    QMessageBox.critical(self, "Ошибка ввода",
                                         "Введено не число")
                    value = None
                else:
                    value = new_value
        elif isinstance(el.value, str):
            value, ok = QInputDialog.getText(self,
                                             title,
                                             "Изменить текст:",
                                             text=el.value)
        elif isinstance(el.value, datetime.date):
            dialog = date_dialogs.DateDialog(self, title, el.value)
            ok = bool(dialog.exec())
            value = dialog.dateEdit.date().toPyDate()
        elif isinstance(el.value, datetime.time):
            dialog = date_dialogs.TimeDialog(self, title, el.value)
            ok = bool(dialog.exec())
            value = dialog.timeEdit.time().toPyTime()
        elif isinstance(el.value, datetime.datetime):
            dialog = date_dialogs.DateTimeDialog(self, title, el.value)
            ok = bool(dialog.exec())
            value = dialog.dateTimeEdit.dateTime().toPyDateTime()
        elif isinstance(el.value, datetime.timedelta):
            dialog = date_dialogs.DateTimeDialog(
                self, title,
                datetime.datetime(2000, 1, 1) + el.value,
                "Изменить промежуток (ВНИМАНИЕ! Минимальная дата - 01.01.2000, она будет вычтена при подсчете промежутка)"
            )
            ok = bool(dialog.exec())
            value = dialog.dateTimeEdit.dateTime().toPyDateTime(
            ) - datetime.datetime(2000, 1, 1)
        elif isinstance(el.value, uuid.UUID):
            value, ok = QInputDialog.getText(self,
                                             title,
                                             "Изменить UUID:",
                                             text=str(el.value))
            if ok and value:
                try:
                    value = uuid.UUID(value)
                except ValueError:
                    value = None

        if ok and value:
            id = self.tableWidget.item(row, 0)
            col_name = self.tableWidget.horizontalHeaderItem(col)
            try:
                self.db.update_data(self.currentTable, value, id.text(),
                                    col_name.text().lower())
                self.updateTable()
            except Error as e:
                QMessageBox.critical(self, "Ошибка базы данных", str(e))
