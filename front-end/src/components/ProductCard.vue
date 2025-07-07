<script setup>
import Wrapper from "@/components/Wrapper.vue";
import {Button} from "@/components/ui/button/index.js";
import {Heart, Plus} from 'lucide-vue-next';

const props = defineProps({
  productName: String,
  productPrice: Number,
  productColorsAmount: Number,
  productImage: String,
})
</script>

<template>
  <div class="w-full max-w-[300px]">
    <Wrapper>
      <div class="flex flex-col gap-2 items-center">
        <div class="rounded-sm w-[150px] h-[150px] overflow-hidden">
          <img class="w-full h-full object-cover" :src="productImage" alt="product image">
        </div>
        <div class="w-full flex flex-col">
          <p class="text-base">{{ productName }}</p>
          <div class="flex justify-between items-center">
            <p class="font-bold">{{ productPrice }} &#8364;</p>
            <div class="p-1 px-2 bg-gray-100 rounded-sm">
              <p>{{ productColorsAmount }} colors</p>
            </div>
          </div>
        </div>
        <div class="flex items-center w-full">
          <Button class="flex-1 bg-[#c9a275] hover:bg-[#dbb384] cursor-pointer">Add to cart
            <Plus></Plus>
          </Button>
          <Heart class="m-5 mr-0 w-6 h-6 cursor-pointer"></Heart>
        </div>
      </div>
    </Wrapper>
  </div>
</template>