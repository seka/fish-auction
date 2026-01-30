import { render, screen } from "@testing-library/react";
import { MainLayoutTemplate } from "./MainLayoutTemplate";
import { describe, it, expect } from "vitest";
import { ToastProvider } from "@/src/hooks/useToast";

describe("MainLayoutTemplate", () => {
    it("renders children correctly", () => {
        render(
            <ToastProvider>
                <MainLayoutTemplate navbar={<div>Navbar Helper</div>}>
                    <div data-testid="main-content">Main Content</div>
                </MainLayoutTemplate>
            </ToastProvider>
        );

        expect(screen.getByTestId("main-content")).toBeInTheDocument();
        expect(screen.getByText("Main Content")).toBeInTheDocument();
    });

    it("renders the navbar correctly", () => {
        render(
            <ToastProvider>
                <MainLayoutTemplate navbar={<div data-testid="navbar">Mock Navbar</div>}>
                    <div>Content</div>
                </MainLayoutTemplate>
            </ToastProvider>
        );

        expect(screen.getByTestId("navbar")).toBeInTheDocument();
        expect(screen.getByText("Mock Navbar")).toBeInTheDocument();
    });
});
