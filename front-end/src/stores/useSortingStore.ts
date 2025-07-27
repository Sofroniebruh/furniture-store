import {defineStore} from "pinia";

export const useSortingStore = defineStore("sortingStore", {
    state: () => ({
        models: [] as string[],
        sorting: "",
        priceRange: [0, 1500] as number[],
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
        addPriceRange(priceFrom: number, priceTo: number) {
            this.priceRange[0] = priceFrom;
            this.priceRange[1] = priceTo;
        }
    }
})