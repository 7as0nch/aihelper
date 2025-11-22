<script setup lang="ts">
import { DatePicker } from 'ant-design-vue';
import type { Dayjs } from 'dayjs';
import type { MentionType, MentionOption } from '../../../config/mentions';

const ARangePicker = DatePicker.RangePicker;

defineProps<{
  show: boolean;
  position: { top: number; left: number };
  options: MentionOption[];
  mentionTypes: MentionType[];
  activeIndex: number;
  activeType: MentionType | null;
  showDatePicker: boolean;
  dateRangeValue: [Dayjs, Dayjs] | undefined;
  pickerMode: 'date' | 'time' | 'datetime';
  getDateFormat: string;
}>();

const emit = defineEmits<{
  (e: 'select-type', type: MentionType): void;
  (e: 'select-option', option: MentionOption): void;
  (e: 'close'): void;
  (e: 'update:dateRangeValue', value: [Dayjs, Dayjs] | undefined): void;
  (e: 'update:pickerMode', mode: 'date' | 'time' | 'datetime'): void;
  (e: 'confirm-date'): void;
  (e: 'update:activeType', type: MentionType | null): void;
  (e: 'update:showDatePicker', show: boolean): void;
}>();

const handleDateRangeChange = (val: any) => {
  emit('update:dateRangeValue', val as [Dayjs, Dayjs] | undefined);
};
</script>

<template>
  <div 
    v-if="show"
    class="absolute bottom-full left-0 sm:left-4 mb-2 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-20 mention-menu"
    :class="showDatePicker ? 'w-full sm:w-96' : 'w-64'"
    @mousedown.stop
  >
    <div class="p-1">
      <!-- Date Range Picker State -->
      <template v-if="showDatePicker">
         <div class="px-2 py-1 text-xs text-gray-400 font-medium flex items-center gap-1 mb-2">
           <button @mousedown.prevent="emit('update:showDatePicker', false)" class="hover:text-primary">返回</button>
           <span>/</span>
           <span>自定义范围</span>
         </div>
         
         <div class="p-2 space-y-3" @mousedown.stop>
           <!-- Mode Selection -->
           <div class="flex items-center gap-2 mb-2">
             <label class="text-xs text-gray-500">模式:</label>
             <div class="flex bg-gray-100 dark:bg-gray-800 rounded-lg p-0.5">
               <button 
                 v-for="mode in (['date', 'time', 'datetime'] as const)" 
                 :key="mode"
                 @mousedown.prevent="emit('update:pickerMode', mode)"
                 class="px-2 py-1 text-xs rounded-md transition-colors capitalize"
                 :class="pickerMode === mode ? 'bg-white dark:bg-gray-700 shadow-sm text-primary' : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'"
               >
                 {{ mode === 'datetime' ? 'Date+Time' : mode }}
               </button>
             </div>
           </div>

           <!-- Ant Design Vue Range Picker -->
             <div class="w-full">
               <ARangePicker 
                 :value="dateRangeValue"
                 @update:value="handleDateRangeChange"
                 :show-time="pickerMode === 'datetime' || pickerMode === 'time'"
                 :picker="pickerMode === 'time' ? 'time' : 'date'"
                 :format="getDateFormat"
                 class="w-full"
                 placement="topLeft"
                 dropdown-class-name="chat-date-picker-dropdown"
               />
             </div>
           
           <button 
             @mousedown.prevent="emit('confirm-date')"
             class="w-full mt-3 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
             :disabled="!dateRangeValue"
           >
             确认引用
           </button>
         </div>
      </template>

      <template v-else-if="!activeType">
        <div class="px-2 py-1 text-xs text-gray-400 font-medium">选择类型</div>
        <button 
          v-for="(type, index) in mentionTypes" 
          :key="type.key"
          @mousedown.prevent="emit('select-type', type)"
          class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
          :class="{ 'bg-blue-50 dark:bg-blue-900/20': index === activeIndex }"
        >
          <div class="w-6 h-6 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
            <component :is="type.icon" class="w-3 h-3 text-gray-600 dark:text-gray-300" />
          </div>
          <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ type.label }}</span>
        </button>
      </template>
      
      <template v-else>
         <div class="px-2 py-1 text-xs text-gray-400 font-medium flex items-center gap-1">
           <button @mousedown.prevent="emit('update:activeType', null)" class="hover:text-primary">返回</button>
           <span>/</span>
           <span>{{ activeType.label }}</span>
         </div>
         <button 
          v-for="(option, index) in options" 
          :key="option.id"
          @mousedown.prevent="emit('select-option', option)"
          class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
          :class="{ 'bg-blue-50 dark:bg-blue-900/20': index === activeIndex }"
        >
          <div class="flex flex-col">
            <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ option.label }}</span>
            <span v-if="option.description" class="text-xs text-gray-500">{{ option.description }}</span>
          </div>
        </button>
        <div v-if="options.length === 0" class="px-3 py-4 text-center text-gray-500 text-sm">
          无匹配结果
        </div>
      </template>
    </div>
  </div>
</template>
