export {};

declare global {
  // Attaching exports to the global scope in designated namespace
  interface Window {
    _dt: Object;
    // TODO: Write Date definition in typescript
    dates: Object[];
  }
  interface String {
    capitalize(): string;
  }
}