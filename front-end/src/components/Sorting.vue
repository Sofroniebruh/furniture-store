<script setup>

import {Button} from "@/components/ui/button/index.js";
import {useParams} from "@/composables/useParams.ts";
import {computed, watch} from "vue";
import {Check} from 'lucide-vue-next';
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";
import {useSortingStore} from "@/stores/useSortingStore.js";
import PriceRange from "@/components/PriceRange.vue";

const sortingStore = useSortingStore();

const priceRange = computed(() => sortingStore.priceRange)
const models = computed(() => sortingStore.models);
const sortingType = computed(() => sortingStore.sorting)

const {updateModel, updateSorting, updatePriceRange} = useParams()
const {setAllSheetsClosed} = useScreenSheetStore()

if (window.innerWidth >= 768) {
  setAllSheetsClosed()
}

const isCheckedTable = computed(() => models.value.some((m) => m === 'table'))
const isCheckedSofa = computed(() => models.value.some((m) => m === 'sofa'))
const isCheckedBed = computed(() => models.value.some((m) => m === 'bed'))
const isCheckedChair = computed(() => models.value.some((m) => m === 'chair'))

const handleNewModel = (model) => {
  const isPresent = models.value.some((m) => m === model)

  if (isPresent) {
    sortingStore.removeModel(model)
  } else {
    sortingStore.addModel(model)
  }
}

const handleSearch = async () => {
  setAllSheetsClosed()

  if (window.innerWidth < 768) {
    await updateModel(models.value)
    await updateSorting(sortingType.value)
    await updatePriceRange(priceRange.value[0], priceRange.value[1])
  }
}

watch(models, (newValue) => {
  if (window.innerWidth >= 768) {
    updateModel(newValue)
  }
})
watch(sortingType, (newValue) => {
  if (window.innerWidth >= 768) {
    updateSorting(newValue)
  }
})
watch(priceRange, (newValue) => {
  if (window.innerWidth >= 768) {
    updatePriceRange(newValue)
  }
})
</script>

<template>
  <div class="sticky top-[88px]">
    <div class="flex flex-col gap-3 mb-2">
      <h1 class="text-base font-semibold">Sort by</h1>
      <div class="flex flex-col gap-3 mb-2">
        <div class="space-x-1.5" v-for="(sorting, index) in ['lth', 'htl', 'sale']">
          <input @click="sortingStore.addSorting(sorting)" :checked="sortingType === sorting" :key="index" type="radio"
                 :name="sorting" :id="sorting" :value="sorting">
          <label v-if="sorting === 'lth'" :key="index" :for="sorting">Price: Low to High</label>
          <label v-if="sorting === 'htl'" :key="index" :for="sorting">Price: High to Low</label>
          <label v-if="sorting === 'sale'" :key="index" :for="sorting">Sale</label>
        </div>
      </div>
    </div>
    <div class="flex flex-col gap-3 mb-2">
      <h1 class="text-base font-semibold">Price range</h1>
      <PriceRange></PriceRange>
    </div>
    <div class="flex flex-col gap-3">
      <h1 class="text-base font-semibold">Model</h1>
      <div class="flex gap-1 items-center leading-0">
        <label class="relative flex items-center cursor-pointer">
          <input
              name="table"
              type="checkbox"
              :checked="isCheckedChair"
              @change="handleNewModel('chair')"
              class="peer hidden w-4 h-4"
          />
          <span
              class="w-4 h-4 flex items-center justify-center border border-gray-200 rounded-[4px]
             peer-checked:bg-[#c9a275] peer-checked:border-transparent"
          >
            <span class="text-white text-sm block"><Check class="w-3 h-3"/></span>
          </span>
        </label>
        <p class="text-gray-600">Chair</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <label class="relative flex items-center cursor-pointer">
          <input
              name="table"
              type="checkbox"
              :checked="isCheckedBed"
              @change="handleNewModel('bed')"
              class="peer hidden w-4 h-4"
          />
          <span
              class="w-4 h-4 flex items-center justify-center border border-gray-200 rounded-[4px]
             peer-checked:bg-[#c9a275] peer-checked:border-transparent"
          >
            <span class="text-white text-sm block"><Check class="w-3 h-3"/></span>
          </span>
        </label>
        <p class="text-gray-600">Bed</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <label class="relative flex items-center cursor-pointer">
          <input
              name="table"
              type="checkbox"
              :checked="isCheckedTable"
              @change="handleNewModel('table')"
              class="peer hidden w-4 h-4"
          />
          <span
              class="w-4 h-4 flex items-center justify-center border border-gray-200 rounded-[4px]
             peer-checked:bg-[#c9a275] peer-checked:border-transparent"
          >
            <span class="text-white text-sm block"><Check class="w-3 h-3"/></span>
          </span>
        </label>
        <p class="text-gray-600">Table</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <label class="relative flex items-center cursor-pointer">
          <input
              name="table"
              type="checkbox"
              :checked="isCheckedSofa"
              @change="handleNewModel('sofa')"
              class="peer hidden w-4 h-4"
          />
          <span
              class="w-4 h-4 flex items-center justify-center border border-gray-200 rounded-[4px]
             peer-checked:bg-[#c9a275] peer-checked:border-transparent"
          >
            <span class="text-white text-sm block"><Check class="w-3 h-3"/></span>
          </span>
        </label>
        <p class="text-gray-600">Sofa</p>
      </div>
    </div>
    <Button @click="handleSearch" class="mt-4 w-full cursor-pointer block md:hidden max-w-[200px]">Search
    </Button>
  </div>
</template>