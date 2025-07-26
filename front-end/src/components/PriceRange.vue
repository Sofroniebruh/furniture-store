<template>
  <div class="space-y-4">
    <div class="flex gap-2 max-w-[200px]">
      <input
          name="priceFrom"
          class="appearance-none [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-inner-spin-button]:m-0 text-gray-600 w-full border rounded-sm p-1 text-sm px-2"
          placeholder="From..."
          type="number"
          :max="1500"
          :min="0"
          v-model.number="priceFrom"
          @input="updateFromInput"
          @blur="validateAndUpdate"
      >
      <input
          name="priceTo"
          class="appearance-none [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-inner-spin-button]:m-0 text-gray-600 w-full border rounded-sm p-1 px-2 text-sm"
          placeholder="To..."
          type="number"
          :max="1500"
          :min="0"
          v-model.number="priceTo"
          @input="updateToInput"
          @blur="validateAndUpdate"
      >
    </div>
    <div class="max-w-[200px]">
      <RangeSlider
          :min="0"
          :max="1500"
          :step="1"
          :model-value="priceRange"
          :on-value-change="updatePricesFromSlider"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, nextTick, ref, watch} from 'vue'
import {RangeSlider} from "@/components/ui/range-slider";
import {useSortingStore} from "@/stores/useSortingStore";
import {useParams} from "@/composables/useParams";

const sortingStore = useSortingStore()

const priceFrom = ref<number>(sortingStore.priceRange[0])
const priceTo = ref<number>(sortingStore.priceRange[1])

const priceRange = computed(() => [priceFrom.value || 0, priceTo.value || 1500])
const {updatePriceRange} = useParams()

let updateTimeout: ReturnType<typeof setTimeout> | null = null

const debouncedUpdate = () => {
  if (updateTimeout) {
    clearTimeout(updateTimeout)
  }
  updateTimeout = setTimeout(() => {
    sortingStore.addPriceRange(priceRange.value[0], priceRange.value[1])
    updatePriceRange(priceRange.value[0], priceRange.value[1])
  }, 300)
}

const updatePricesFromSlider = (values: number[]) => {
  priceFrom.value = values[0]
  priceTo.value = values[1]
}

watch(priceRange, () => {
  if (window.innerWidth >= 768) {
    debouncedUpdate()
  } else {
    sortingStore.addPriceRange(priceRange.value[0], priceRange.value[1])
  }
})

const validateRange = () => {
  if (isNaN(priceFrom.value) || priceFrom.value < 0) {
    priceFrom.value = 0
  }
  if (isNaN(priceTo.value) || priceTo.value > 1500) {
    priceTo.value = 1500
  }

  if (priceFrom.value > priceTo.value) {
    priceFrom.value = priceTo.value
  }
}

const updateFromInput = () => {
  nextTick(() => {
    if (priceFrom.value > priceTo.value) {
      priceFrom.value = priceTo.value
    }
    if (priceFrom.value < 0) {
      priceFrom.value = 0
    }
  })
}

const updateToInput = () => {
  nextTick(() => {
    if (priceTo.value < priceFrom.value) {
      priceTo.value = priceFrom.value
    }
    if (priceTo.value > 1500) {
      priceTo.value = 1500
    }
  })
}

const validateAndUpdate = () => {
  validateRange()
  if (updateTimeout) {
    clearTimeout(updateTimeout)
  }
  sortingStore.addPriceRange(priceRange.value[0], priceRange.value[1])
  updatePriceRange(priceRange.value[0], priceRange.value[1])
}

watch(() => sortingStore.priceRange, (newRange) => {
  if (newRange[0] !== priceFrom.value || newRange[1] !== priceTo.value) {
    priceFrom.value = newRange[0]
    priceTo.value = newRange[1]
  }
}, {deep: true})
</script>