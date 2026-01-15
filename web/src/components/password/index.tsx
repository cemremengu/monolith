import { EyeIcon, EyeOffIcon } from "lucide-react";
import * as React from "react";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { cn } from "@/lib/utils";

type PasswordProps = Omit<React.ComponentProps<"input">, "type"> & {
  showStrengthIndicator?: boolean;
};

type PasswordStrength = {
  score: number;
  label: string;
};

function calculatePasswordStrength(password: string): PasswordStrength {
  if (!password) {
    return { score: 0, label: "Very Weak" };
  }

  let score = 0;

  if (password.length >= 8) score++;
  if (password.length >= 12) score++;

  if (/[a-z]/.test(password)) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/\d/.test(password)) score++;
  if (/[^a-zA-Z0-9]/.test(password)) score++;

  const normalizedScore = Math.min(Math.floor((score / 6) * 5), 4);

  const labels = ["Very Weak", "Weak", "Fair", "Strong", "Very Strong"];

  return {
    score: normalizedScore,
    label: labels[normalizedScore],
  };
}

function Password({
  className,
  showStrengthIndicator = false,
  value,
  onChange,
  ...props
}: PasswordProps) {
  const [showPassword, setShowPassword] = React.useState(false);
  const [internalValue, setInternalValue] = React.useState("");

  const currentValue = value !== undefined ? String(value) : internalValue;

  const strength = showStrengthIndicator
    ? calculatePasswordStrength(currentValue)
    : null;

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (value === undefined) {
      setInternalValue(e.target.value);
    }
    onChange?.(e);
  };

  const strengthColors = [
    "bg-red-500",
    "bg-orange-500",
    "bg-yellow-500",
    "bg-lime-500",
    "bg-green-500",
  ];

  const strengthTextColors = [
    "text-red-600",
    "text-orange-600",
    "text-yellow-600",
    "text-lime-600",
    "text-green-600",
  ];

  return (
    <div className="flex flex-col gap-2">
      <InputGroup className={className}>
        <InputGroupInput
          type={showPassword ? "text" : "password"}
          value={value}
          onChange={handleChange}
          {...props}
        />
        <InputGroupAddon align="inline-end">
          <InputGroupButton
            title="Toggle password visibility"
            size="icon-xs"
            onClick={() => setShowPassword((prev) => !prev)}
          >
            {showPassword ? (
              <EyeOffIcon className="size-5" />
            ) : (
              <EyeIcon className="size-5" />
            )}
          </InputGroupButton>
        </InputGroupAddon>
      </InputGroup>
      {showStrengthIndicator && currentValue && strength && (
        <div className="space-y-1">
          <div className="flex gap-1">
            {[0, 1, 2, 3, 4].map((index) => (
              <div
                key={index}
                className={cn(
                  "h-1 flex-1 rounded-full transition-colors",
                  index <= strength.score
                    ? strengthColors[strength.score]
                    : "bg-neutral-200 dark:bg-neutral-700",
                )}
              />
            ))}
          </div>
          <p
            className={cn(
              "text-xs font-medium",
              strengthTextColors[strength.score],
            )}
          >
            {strength.label}
          </p>
        </div>
      )}
    </div>
  );
}

export { Password };
