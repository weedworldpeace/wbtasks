import subprocess
import time

def run_command(cmd, input_file):
    start = time.time()
    with open(input_file, 'r') as f:
        result = subprocess.run(cmd, stdin=f, capture_output=True, text=True, shell=True)
    end = time.time()
    return end - start

def main():
    input_file = "input.txt"
    
    cut_time = run_command("cut -f 1 -d ' '", input_file)
    print(f"cut time: {cut_time:.4f} seconds")
    
    main_time = run_command("./main -f 1 -d ' '", input_file)
    print(f"./main time: {main_time:.4f} seconds", )

if __name__ == "__main__":
    main()