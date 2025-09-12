import {defineStore} from "pinia";

type ScreenSheet = {
    name: "SmallScreenMenu" | "SmallScreenFiltering" | "SmallScreenSearch" | "ScreenCart" | "ScreenProfile";
}

type ScreenDialog = {
    name: "RegisterDialog" | "LoginDialog" | "CreationDialog" | "EditDialog" | "DeleteDialog";
}

export const useScreenSheetStore = defineStore("smallScreenSheetStore", {
        state: () => ({
            sheets: [] as ScreenSheet[],
            dialogs: [] as ScreenDialog[],
        }),
        actions: {
            setOpenSheet(open: boolean, sheet: ScreenSheet) {
                let updatedSheets = this.sheets

                if (open) {
                    const exists = updatedSheets.some((s: ScreenSheet) => s.name === sheet.name)
                    updatedSheets = exists ? updatedSheets : [...updatedSheets, sheet]
                } else {
                    updatedSheets = updatedSheets.filter((s: ScreenSheet) => s.name !== sheet.name)
                }

                this.sheets = updatedSheets
            },
            setOpenDialog(open: boolean, sheet: ScreenDialog) {
                let updatedDialogs = this.dialogs

                if (open) {
                    const exists = updatedDialogs.some((s: ScreenDialog) => s.name === sheet.name)
                    updatedDialogs = exists ? updatedDialogs : [...updatedDialogs, sheet]
                } else {
                    updatedDialogs = updatedDialogs.filter((s: ScreenDialog) => s.name !== sheet.name)
                }

                this.dialogs = updatedDialogs
            },
            isDialogOpen(name: string) {
                return this.dialogs.some((s: ScreenDialog) => s.name === name)
            },
            isSheetOpen(name: string) {
                return this.sheets.some((s: ScreenSheet) => s.name === name)
            },
            setAllSheetsClosed() {
                this.sheets = []
            },
            setAllDialogsClosed() {
                this.dialogs = []
            }
        }
    }
)