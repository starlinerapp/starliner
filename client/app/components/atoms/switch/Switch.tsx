import React from "react";
import * as SwitchPrimitive from "@radix-ui/react-switch";

type SwitchProps = React.ComponentPropsWithoutRef<typeof SwitchPrimitive.Root>;

const Switch = React.forwardRef<
  React.ComponentRef<typeof SwitchPrimitive.Root>,
  SwitchProps
>((props, ref) => {
  return (
    <SwitchPrimitive.Root
      ref={ref}
      className="bg-black-a4 data-[state=checked]:bg-violet-10 relative h-[25px] w-[42px] cursor-default cursor-pointer rounded-full outline-none"
      {...props}
    >
      <SwitchPrimitive.Thumb className="shadow-black-a4 block size-[21px] translate-x-0.5 rounded-full bg-white shadow-[0_2px_2px] transition-transform duration-100 will-change-transform data-[state=checked]:translate-x-[19px]" />
    </SwitchPrimitive.Root>
  );
});

Switch.displayName = "Switch";

export default Switch;
