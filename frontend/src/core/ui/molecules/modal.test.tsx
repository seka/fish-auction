import { render, screen } from "@testing-library/react";
import { ModalBackdrop, ModalContent } from "./modal";
import { describe, it, expect } from "vitest";

describe("Modal", () => {
    describe("ModalBackdrop", () => {
        it("renders children correctly", () => {
            render(
                <ModalBackdrop>
                    <div data-testid="backdrop-content">Backdrop Content</div>
                </ModalBackdrop>
            );
            expect(screen.getByTestId("backdrop-content")).toBeInTheDocument();
        });
    });

    describe("ModalContent", () => {
        it("renders children correctly", () => {
            render(
                <ModalContent>
                    <div data-testid="modal-content">Modal Content</div>
                </ModalContent>
            );
            expect(screen.getByTestId("modal-content")).toBeInTheDocument();
        });
    });
});
