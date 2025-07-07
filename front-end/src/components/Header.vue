<script setup>
import Wrapper from "@/components/Wrapper.vue";
import {ArrowRight, Heart, Menu, ShoppingBag, User} from 'lucide-vue-next';
import SheetGeneral from "@/components/SheetGeneral.vue";
import {useRoute} from "vue-router";
import {computed} from "vue";
import {cn} from "@/lib/utils.js";

const route = useRoute()
const isOnAboutPage = computed(() => route.path === "/about");
const isOnContactPage = computed(() => route.path === "/contact");
const isOnProductsPage = computed(() => route.path === "/products");
</script>

<template>
  <div class="z-50 w-full backdrop-blur-sm bg-white/70 fixed top-0">
    <Wrapper>
      <div class="flex items-center justify-between">
        <SheetGeneral trigger-class="sm:hidden" side="left" title="FUMI">
          <template #trigger>
            <Menu class="order-1 block sm:hidden cursor-pointer"/>
          </template>
          <template #content>
            <ul class="p-4">
              <li>
                <div class=" text-gray-700 hover:text-black">
                  <router-link to="/products" class="text-2xl flex justify-between w-full items-center">
                    <p class="text-2xl ">Products</p>
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
              <li>
                <div class="text-gray-700 hover:text-black">
                  <router-link to="/about" class="text-2xl flex justify-between w-full items-center">
                    About Us
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
              <li>
                <div class="text-gray-700 hover:text-black">
                  <router-link to="/contact" class="text-2xl flex w-full justify-between items-center">
                    Contact Us
                    <ArrowRight/>
                  </router-link>
                </div>
              </li>
            </ul>
            <ul class="p-4">
              <li class="flex items-center gap-2 text-2xl text-gray-700 hover:text-black">
                <User/>
                Profile
              </li>
              <li class="flex items-center gap-2 text-2xl text-gray-700 hover:text-black">
                <Heart/>
                Wishlist
              </li>
            </ul>
          </template>
        </SheetGeneral>
        <router-link to="/" class="order-2 sm:order-1">
          <h1 class="font-semibold text-base sm:text-lg">FUMI</h1>
        </router-link>
        <ul class="sm:order-2 hidden sm:flex gap-3">
          <li :class="cn('text-gray-700 hover:text-black', isOnProductsPage && 'text-black')">
            <router-link to="/products">
              Products
            </router-link>
          </li>
          <li :class="cn('text-gray-700 hover:text-black', isOnAboutPage && 'text-black')">
            <router-link to="/about">
              About Us
            </router-link>
          </li>
          <li :class="cn('text-gray-700 hover:text-black', isOnContactPage && 'text-black')">
            <router-link to="/contact">
              Contact Us
            </router-link>
          </li>
        </ul>
        <ul class="order-3 flex gap-3">
          <li class="hidden sm:block">
            <User></User>
          </li>
          <li class="hidden sm:block">
            <Heart></Heart>
          </li>
          <SheetGeneral title="Your Cart">
            <template #trigger>
              <li>
                <ShoppingBag></ShoppingBag>
              </li>
            </template>
            <template #content>
              <h1>Content!</h1>
            </template>
          </SheetGeneral>
        </ul>
      </div>
    </Wrapper>
  </div>
</template>