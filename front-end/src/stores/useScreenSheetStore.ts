import {defineStore} from "pinia";

type ScreenSheet = {
    name: "SmallScreenMenu" | "SmallScreenFiltering" | "SmallScreenSearch" | "ScreenCart";
}

export const useScreenSheetStore = defineStore("smallScreenSheetStore", {
        state: () => ({
            sheets: [] as ScreenSheet[],
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
            isSheetOpen(name: string) {
                return this.sheets.some((s: ScreenSheet) => s.name === name)
            },
            setAllSheetsClosed() {
                this.sheets = []
            }
        }
    }
)