<script setup>
import {Facebook, Instagram, Twitter} from 'lucide-vue-next';
</script>

<template>
  <footer class="rounded-2xl py-[34px] sm:py-[68px] px-8 flex justify-around items-center flex-col sm:flex-row">
    <ul class="flex flex-col items-center sm:items-start">
      <li class="mt-2">
        <a href="#" class="font-medium text-sm">Privacy policy</a>
      </li>
      <li class="mb-2">
        <a href="#" class="font-medium text-sm">Terms and conditions</a>
      </li>
    </ul>
    <h2 class="font-medium my-4 sm:m-0">&copy; 2024 All rights reserved.</h2>
    <div>
      <ul class="flex gap-2 my-4 sm:m-0">
        <li>
          <a href="#">
            <Facebook/>
          </a>
        </li>
        <li>
          <a href="#">
            <Instagram></Instagram>
          </a>
        </li>
        <li>
          <a href="#">
            <Twitter/>
          </a>
        </li>
      </ul>
    </div>
  </footer>
</template>