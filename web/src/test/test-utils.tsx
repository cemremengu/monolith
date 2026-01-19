import type { ReactNode } from "react";
import { render, type RenderOptions } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { I18nextProvider } from "react-i18next";

import i18n from "@/i18n";
import { ThemeProvider } from "@/hooks/use-theme";

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

type AllProvidersProps = {
  children: ReactNode;
};

// eslint-disable-next-line react-refresh/only-export-components
function AllProviders({ children }: AllProvidersProps) {
  const queryClient = createTestQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <I18nextProvider i18n={i18n}>
        <ThemeProvider defaultTheme="system" storageKey="test-theme">
          {children}
        </ThemeProvider>
      </I18nextProvider>
    </QueryClientProvider>
  );
}

function customRender(ui: React.ReactElement, options?: RenderOptions) {
  return render(ui, { wrapper: AllProviders, ...options });
}

// Override render with our custom render
export { customRender as render };
