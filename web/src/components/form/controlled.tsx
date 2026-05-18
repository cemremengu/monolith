import type * as React from "react";

import { type Control, type FieldPath, type FieldValues, useController } from "react-hook-form";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";

type ControlledFieldProps<T extends FieldValues> = {
  control: Control<T>;
  name: FieldPath<T>;
  label: string;
  id?: string;
};

export function FormInput<T extends FieldValues>({
  control,
  name,
  label,
  id,
  ...inputProps
}: ControlledFieldProps<T> & Omit<React.ComponentProps<typeof Input>, "name">) {
  const { field, fieldState } = useController({ control, name });
  const fieldId = id ?? field.name;

  return (
    <Field data-invalid={fieldState.invalid}>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <Input {...field} {...inputProps} id={fieldId} aria-invalid={fieldState.invalid} />
      {fieldState.invalid && <FieldError>{fieldState.error?.message}</FieldError>}
    </Field>
  );
}

export function FormTextarea<T extends FieldValues>({
  control,
  name,
  label,
  id,
  ...textareaProps
}: ControlledFieldProps<T> & Omit<React.ComponentProps<typeof Textarea>, "name">) {
  const { field, fieldState } = useController({ control, name });
  const fieldId = id ?? field.name;

  return (
    <Field data-invalid={fieldState.invalid}>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <Textarea {...field} {...textareaProps} id={fieldId} aria-invalid={fieldState.invalid} />
      {fieldState.invalid && <FieldError>{fieldState.error?.message}</FieldError>}
    </Field>
  );
}

export function FormSelect<T extends FieldValues>({
  control,
  name,
  label,
  id,
  placeholder,
  children,
}: ControlledFieldProps<T> & {
  placeholder?: string;
  children: React.ReactNode;
}) {
  const { field, fieldState } = useController({ control, name });
  const fieldId = id ?? field.name;

  return (
    <Field data-invalid={fieldState.invalid}>
      <FieldLabel htmlFor={fieldId}>{label}</FieldLabel>
      <Select name={field.name} value={field.value} onValueChange={field.onChange}>
        <SelectTrigger id={fieldId} aria-invalid={fieldState.invalid}>
          <SelectValue placeholder={placeholder} />
        </SelectTrigger>
        <SelectContent>{children}</SelectContent>
      </Select>
      {fieldState.invalid && <FieldError>{fieldState.error?.message}</FieldError>}
    </Field>
  );
}
