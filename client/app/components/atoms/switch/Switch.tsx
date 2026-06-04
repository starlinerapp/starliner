import * as SwitchPrimitive from "@radix-ui/react-switch";
import React from "react";

type SwitchProps = React.ComponentPropsWithoutRef<typeof SwitchPrimitive.Root>;

const Switch = React.forwardRef<
  React.ComponentRef<typeof SwitchPrimitive.Root>,
  SwitchProps
>((props, ref) => {
  return (
    <SwitchPrimitive.Root
      ref={ref}
      className="relative h-[25px] w-[42px] cursor-default cursor-pointer rounded-full bg-black-a4 outline-none data-[state=checked]:bg-violet-10"
      {...props}
    >
      <SwitchPrimitive.Thumb className="block size-[21px] translate-x-0.5 rounded-full bg-white shadow-[0_2px_2px] shadow-black-a4 transition-transform duration-100 will-change-transform data-[state=checked]:translate-x-[19px]" />
    </SwitchPrimitive.Root>
  );
});

Switch.displayName = "Switch";

export default Switch;
