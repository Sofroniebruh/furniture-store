<script setup>

import {RadioGroup, RadioGroupItem} from '@/components/ui/radio-group'
import {RangeSlider} from "@/components/ui/range-slider/index.js";
import {Button} from "@/components/ui/button/index.js";
import {useParams} from "@/composables/useParams.js";
import {computed, ref, watch} from "vue";
import {Check} from 'lucide-vue-next';

let models = ref([])
const {updateModel, updateEvent, updatePriceTo, updatePriceFrom} = useParams()

const isCheckedTable = computed(() => models.value.some((m) => m === 'table'))
const isCheckedSofa = computed(() => models.value.some((m) => m === 'sofa'))
const isCheckedBed = computed(() => models.value.some((m) => m === 'bed'))
const isCheckedChair = computed(() => models.value.some((m) => m === 'chair'))

const handleNewModel = (model) => {
  const isPresent = models.value.some((m) => m === model)

  if (isPresent) {
    models.value = models.value.filter((m) => m !== model)
  } else {
    models.value = [...models.value, model]
  }
}

const handleSearch = () => {
  updateModel(models.value)
}

watch(models, (newValue) => {
  if (window.screen.width >= 768) {
    updateModel(newValue)
  }
}, {immediate: true})

const updatePrices = (prices) => {
  console.log("Prices", prices)
}
</script>

<template>
  <div class="sticky top-[88px]">
    <div class="flex flex-col gap-3 mb-2">
      <h1 class="text-base font-semibold">Sort by</h1>
      <RadioGroup default-value="option-one">
        <div class="flex items-center space-x-1.5">
          <RadioGroupItem id="l-t-h" value="l-t-h"/>
          <Label class="text-gray-600" for="l-t-h">Price: low to high</Label>
        </div>
        <div class="flex items-center space-x-1.5">
          <RadioGroupItem id="h-t-l" value="h-t-l"/>
          <Label class="text-gray-600" for="h-t-l">Price: high to low</Label>
        </div>
        <div class="flex items-center space-x-1.5">
          <RadioGroupItem id="on-sale" value="on-sale"/>
          <Label class="text-gray-600" for="on-sale">On sale</Label>
        </div>
      </RadioGroup>
    </div>
    <div class="flex flex-col gap-3 mb-2">
      <h1 class="text-base font-semibold">Price range</h1>
      <div class="flex flex-col gap-6">
        <div class="flex gap-2 max-w-[200px]">
          <input name="priceFrom"
                 class="appearance-none [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-inner-spin-button]:m-0 text-gray-600 w-full border rounded-sm p-1 text-sm px-2"
                 placeholder="From..."
                 type="number"
                 :max="1500"
                 :min="0"
          >
          <input name="priceTo"
                 class="appearance-none [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-inner-spin-button]:m-0 text-gray-600 w-full border rounded-sm p-1 px-2 text-sm"
                 placeholder="To..."
                 type="number"
                 :max="1500"
                 :min="0"
          >
        </div>
        <div class="max-w-[200px]">
          <RangeSlider :min="0" :max="1500" :step="1" :model-value="[0, 1500]"
                       :on-value-change="updatePrices"></RangeSlider>
        </div>
      </div>
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