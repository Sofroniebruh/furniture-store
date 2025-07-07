<script setup>
import Wrapper from "@/components/Wrapper.vue";
import {ArrowRight, Heart, Menu, Search, ShoppingBag, User} from 'lucide-vue-next';
import SheetGeneral from "@/components/SheetGeneral.vue";
import {useRoute} from "vue-router";
import {computed} from "vue";
import {cn} from "@/lib/utils.js";
import {Button} from "@/components/ui/button/index.js";

const route = useRoute()
const isOnAboutPage = computed(() => route.path === "/about");
const isOnContactPage = computed(() => route.path === "/contact");
const isOnProductsPage = computed(() => route.path === "/products");
</script>

<template>
  <div class="z-50 w-full backdrop-blur-sm bg-white/70 fixed top-0">
    <Wrapper>
      <div class="grid lg:grid-cols-[328px_1fr_auto] sm:grid-cols-[132px_1fr_auto] grid-cols-3">
        <SheetGeneral trigger-class="sm:hidden" side="left" title="FUMI">
          <template #trigger>
            <Menu class="order-1 block sm:hidden cursor-pointer col-start-1"/>
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
        <router-link to="/"
                     class="order-2 sm:order-1 col-start-2 sm:col-start-1 flex items-center justify-center sm:justify-start">
          <h1 class="font-semibold text-base sm:text-lg">FUMI</h1>
        </router-link>
        <ul class="sm:order-2 col-start-2 hidden sm:flex sm:items-center sm:justify-center gap-3">
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
        <ul class="order-3 col-start-3 flex items-center justify-end gap-3">
          <li class="relative items-center hidden lg:flex">
            <input name="search" class="pl-7 border rounded-sm px-2 py-1" type="text" placeholder="Search..."/>
            <Search class="text-gray-600 absolute top-[5px] left-1"/>
          </li>
          <SheetGeneral side="top">
            <template #trigger>
              <li class="lg:hidden block">
                <Search></Search>
              </li>
            </template>
            <template #content>
              <div class="relative items-center w-full flex p-5 pt-0 gap-3">
                <input name="search" class="pl-7 border w-full rounded-sm px-2 py-1" type="text"
                       placeholder="Search..."/>
                <Search class="text-gray-600 absolute top-[5px] left-6"/>
                <Button class="bg-[#c9a275] h-[34px] cursor-pointer hover:bg-[#dbb384]">Go</Button>
              </div>
            </template>
          </SheetGeneral>
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