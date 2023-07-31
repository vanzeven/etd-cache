def extract_columns(input_file, output_file, op):
    with open(input_file, 'r') as infile:
        with open(output_file, op) as outfile:
            for line in infile:
                columns = line.strip().split()
                if len(columns) >= 6:
                    outfile.write(f"{columns[3]} {columns[5]}\n")


if __name__ == "__main__":
    output_filename = "./homes/homes"  # Replace with the name of the output file

    for i in range(21):
        input_filename = "./homes/homes"  # Replace with the name of your input file
        if i == 0:
            op = 'w'
        else:
            op = 'a'
        if i+1 <= 9:
            input_filename2 = input_filename + "0"
        else:
            input_filename2 = input_filename
        extract_columns(input_filename2 + str(i+1) + ".blkparse", output_filename, op)
        print("file " + str(i+1) + " done")

