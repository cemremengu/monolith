import { useCallback } from "react";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const TIMEZONES = [
  { value: "Pacific/Midway", label: "(UTC-11:00) Midway Island" },
  { value: "Pacific/Honolulu", label: "(UTC-10:00) Hawaii" },
  { value: "Pacific/Marquesas", label: "(UTC-09:30) Marquesas Islands" },
  { value: "America/Anchorage", label: "(UTC-09:00) Alaska" },
  {
    value: "America/Los_Angeles",
    label: "(UTC-08:00) Pacific Time (US & Canada)",
  },
  { value: "America/Denver", label: "(UTC-07:00) Mountain Time (US & Canada)" },
  { value: "America/Chicago", label: "(UTC-06:00) Central Time (US & Canada)" },
  {
    value: "America/New_York",
    label: "(UTC-05:00) Eastern Time (US & Canada)",
  },
  { value: "America/Caracas", label: "(UTC-04:00) Atlantic Time (Canada)" },
  { value: "America/St_Johns", label: "(UTC-03:30) Newfoundland" },
  { value: "America/Sao_Paulo", label: "(UTC-03:00) Brasilia" },
  { value: "Atlantic/South_Georgia", label: "(UTC-02:00) South Georgia" },
  { value: "Atlantic/Azores", label: "(UTC-01:00) Azores" },
  { value: "UTC", label: "(UTC+00:00) Coordinated Universal Time" },
  { value: "Europe/London", label: "(UTC+00:00) Greenwich Mean Time" },
  { value: "Europe/Berlin", label: "(UTC+01:00) Central European Time" },
  { value: "Europe/Helsinki", label: "(UTC+02:00) Eastern European Time" },
  { value: "Europe/Moscow", label: "(UTC+03:00) Moscow Time" },
  { value: "Asia/Dubai", label: "(UTC+04:00) Gulf Standard Time" },
  { value: "Asia/Kabul", label: "(UTC+04:30) Afghanistan Time" },
  { value: "Asia/Karachi", label: "(UTC+05:00) Pakistan Standard Time" },
  { value: "Asia/Kolkata", label: "(UTC+05:30) India Standard Time" },
  { value: "Asia/Kathmandu", label: "(UTC+05:45) Nepal Time" },
  { value: "Asia/Dhaka", label: "(UTC+06:00) Bangladesh Standard Time" },
  { value: "Asia/Yangon", label: "(UTC+06:30) Myanmar Time" },
  { value: "Asia/Bangkok", label: "(UTC+07:00) Indochina Time" },
  { value: "Asia/Shanghai", label: "(UTC+08:00) China Standard Time" },
  { value: "Asia/Pyongyang", label: "(UTC+08:30) Pyongyang Time" },
  { value: "Asia/Tokyo", label: "(UTC+09:00) Japan Standard Time" },
  {
    value: "Australia/Adelaide",
    label: "(UTC+09:30) Australian Central Standard Time",
  },
  {
    value: "Australia/Sydney",
    label: "(UTC+10:00) Australian Eastern Standard Time",
  },
  {
    value: "Australia/Lord_Howe",
    label: "(UTC+10:30) Lord Howe Standard Time",
  },
  { value: "Pacific/Norfolk", label: "(UTC+11:00) Norfolk Time" },
  { value: "Pacific/Auckland", label: "(UTC+12:00) New Zealand Standard Time" },
  { value: "Pacific/Chatham", label: "(UTC+12:45) Chatham Standard Time" },
  { value: "Pacific/Tongatapu", label: "(UTC+13:00) Tonga Time" },
  { value: "Pacific/Kiritimati", label: "(UTC+14:00) Line Islands Time" },
];

type TimezoneSelectorProps = {
  value?: string;
  onChange: (value: string) => void;
};

export function TimezoneSelector({ value, onChange }: TimezoneSelectorProps) {
  const handleValueChange = useCallback(
    (selected: string | null) => {
      if (selected) {
        onChange(selected);
      }
    },
    [onChange],
  );

  return (
    <Select items={TIMEZONES} value={value ?? "UTC"} onValueChange={handleValueChange}>
      <SelectTrigger className="w-72">
        <SelectValue />
      </SelectTrigger>
      <SelectContent align="end">
        {TIMEZONES.map((timezone) => (
          <SelectItem key={timezone.value} value={timezone.value}>
            {timezone.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
