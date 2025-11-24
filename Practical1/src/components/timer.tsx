import { TimerIcon } from "lucide-react";

interface TimerProps {
  timeLeft: number;
}
export default function Timer({ timeLeft }: TimerProps) {
  let timerClass = "flex items-center justify-center space-x-2 text-2xl font-bold text-gray-700 mb-8";
  if (timeLeft <= 10) {
    timerClass += " warning";
  }
  if (timeLeft <= 5) {
    timerClass += " red danger";
  }
  return (
    <div
      className={timerClass}
      data-testid="timer"
    >
      <TimerIcon className="w-6 h-6" />
      <span>{timeLeft}s</span>
    </div>
  );
}
