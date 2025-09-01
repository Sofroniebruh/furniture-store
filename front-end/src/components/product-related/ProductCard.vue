<script setup lang="ts">
import Wrapper from "@/components/Wrapper.vue";
import {Button} from "@/components/ui/button";
import {Plus} from 'lucide-vue-next';
import {cn} from "@/lib/utils.js";
import {Product} from "@/lib/types";
import {PropType} from "vue";

const props = defineProps({
  product: {
    type: Object as PropType<Product>,
    required: true
  },
  class: {
    type: String,
    required: false,
  }
})
const className = props.class
const product = props.product

console.log(product)
</script>

<template>
  <div :class="cn('w-full max-w-[300px]', className)">
    <Wrapper class="p-0 py-5 sm:p-5">
      <div class="flex flex-col gap-2 items-center">
        <div class="rounded-sm w-[150px] h-[150px] overflow-hidden">
          <img class="w-full h-full object-cover" :src="product.pictureUrls[0]" alt="product image">
        </div>
        <div class="w-full flex flex-col">
          <p class="text-base">{{ product.name }}</p>
          <div class="flex justify-between items-center">
            <p class="font-bold">{{ product.price }} &#8364;</p>
            <div class="p-1 px-2 bg-gray-100 rounded-sm">
              <p v-if="product.colors.length > 1">{{ product.colors.length }} colors</p>
              <p v-if="product.colors.length === 1">{{ product.colors.length }} color</p>
            </div>
          </div>
        </div>
        <div class="flex items-center w-full mt-[17px]">
          <Button class="flex-1 bg-[#c9a275] hover:bg-[#dbb384] cursor-pointer">Add to cart
            <Plus></Plus>
          </Button>
        </div>
      </div>
    </Wrapper>
  </div>
</template>