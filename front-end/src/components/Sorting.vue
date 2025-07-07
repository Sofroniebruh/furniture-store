<script setup>

import {Checkbox} from "@/components/ui/checkbox/index.js";
import {RadioGroup, RadioGroupItem} from '@/components/ui/radio-group'
import {RangeSlider} from "@/components/ui/range-slider/index.js";
import {Button} from "@/components/ui/button/index.js";

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
        <Checkbox/>
        <p class="text-gray-600">Chair</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <Checkbox/>
        <p class="text-gray-600">Bed</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <Checkbox/>
        <p class="text-gray-600">Table</p>
      </div>
      <div class="flex gap-1 items-center leading-0">
        <Checkbox/>
        <p class="text-gray-600">Sofa</p>
      </div>
    </div>
    <Button class="mt-4 w-full cursor-pointer block md:hidden max-w-[200px]">Search</Button>
  </div>
</template>