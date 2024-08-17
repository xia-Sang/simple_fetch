import subprocess

def get_system_info():
    try:
        result = subprocess.run(['wmic', 'os', 'get', '/all', '/format:list'], capture_output=True, text=True, check=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        return f"An error occurred: {e}"

def parse_system_info(info):
    # Split the info into lines
    lines = info.split('\n')
    
    # Create a dictionary to store the parsed information
    parsed_info = {}
    
    # Iterate through each line
    for line in lines:
        # Split each line into key and value
        parts = line.split('=', 1)
        if len(parts) == 2:
            key = parts[0].strip()
            value = parts[1].strip()
            parsed_info[key] = value
    
    return parsed_info

# Example usage
if __name__ == "__main__":
    system_info = get_system_info()
    print("Raw System Info:")
    print(system_info)
    print("\nParsed System Info:")
    print(parse_system_info(system_info))
