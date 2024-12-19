import csv
import sys
from os import path

from PyQt6.QtWidgets import QApplication, QFileDialog, QMainWindow, QMessageBox
from ui.lib.project import Project
from ui.ui import Ui_MainWindow


class MainApp(QMainWindow, Ui_MainWindow):

    def __init__(self):
        self.dirs = []
        if path.exists(path.expanduser(path.join("~", "data.csv"))):
            with open(path.expanduser(path.join("~", "data.csv"))) as f:
                dirs = csv.reader(f)
                self.dirs = next(dirs)
        super().__init__()
        self.setupUi(self)
        self.restoreTabs()
        self.chooseProject.clicked.connect(self.selectFolder)
        self.openProject.clicked.connect(self.openFolder)
        self.tabWidget.tabCloseRequested.connect(self.closeTab)

    def closeTab(self, index: int):
        if index == 0:
            QMessageBox.warning(
                self,
                "Inappropriate behavour",
                "Нельзя закрывать эту вкладку!",
                QMessageBox.StandardButton.Ok,
            )
            return
        self.dirs.pop(index - 1)
        self.tabWidget.removeTab(index)

    def restoreTabs(self):
        for dir in self.dirs:
            self.tabWidget.addTab(Project(dir), path.basename(dir))

    def selectFolder(self):
        path = QFileDialog.getExistingDirectory(self, "Выбрать папку...")
        self.path.setText(path)

    def openFolder(self):
        dir = self.path.text()
        self.dirs.append(dir)
        i = self.tabWidget.addTab(Project(dir), path.basename(dir))
        self.tabWidget.setCurrentIndex(i)

    def closeEvent(self, event):
        with open(path.expanduser(path.join("~", "data.csv")), "w") as f:
            w = csv.writer(f)
            w.writerow(self.dirs)


def run():
    app = QApplication(sys.argv)
    ex = MainApp()
    ex.show()
    return app.exec()


if __name__ == "__main__":
    sys.exit(run())
