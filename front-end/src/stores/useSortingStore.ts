import {defineStore} from "pinia";

export const useSortingStore = defineStore("sortingStore", {
    state: () => ({
        models: [] as string[],
        sorting: "",
    }),
    actions: {
        addModel(model: string) {
            if (!this.models.includes(model)) this.models = [...this.models, model];
        },
        removeModel(model: string) {
            this.models = this.models.filter((s) => s !== model);
        },
        addSorting(sortingType: string) {
            this.sorting = sortingType;
        },
    }
})