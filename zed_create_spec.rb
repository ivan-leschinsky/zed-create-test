#!/usr/bin/env ruby
require 'fileutils'

file_path = ARGV[0]

# Exit if file doesn't have .rb extension
unless file_path.end_with?('.rb')
  puts "File must have .rb extension"
  exit
end

# Split the path into segments
path_segments = file_path.split('/')

# Check if the file is a spec file or a code file
if file_path.include?('_spec.rb') && path_segments[0] == 'spec'
  # It's a spec file, find the corresponding code file
  path_segments[0] = 'app'
  code_file_path = path_segments.join('/')
  code_file_path = code_file_path.gsub('_spec.rb', '.rb')
  system("zed #{code_file_path}")
elsif path_segments[0] == 'app'
  # It's a code file, find or create the corresponding spec file
  path_segments[0] = 'spec'
  spec_file_path = path_segments.join('/')
  spec_file_path = spec_file_path.gsub('.rb', '_spec.rb')

  # Create the spec file if it doesn't exist
  unless File.exist?(spec_file_path)
    # Ensure the directory exists
    spec_dir = File.dirname(spec_file_path)
    FileUtils.mkdir_p(spec_dir) unless Dir.exist?(spec_dir)

    # Create an empty spec file
    File.open(spec_file_path, 'w') do |file|
      file.puts "# frozen_string_literal: true"
      file.puts
      file.puts "require \"rails_helper\""
      file.puts

      # Extract class name from the source file
      class_name = nil
      if File.exist?(file_path)
        File.foreach(file_path) do |line|
          if line.strip.start_with?("class ")
            # Extract the class name, handling potential inheritance and namespaces
            match = line.match(/class\s+([A-Z][A-Za-z0-9_:]*)\b/)
            class_name = match[1] if match
            break
          end
        end
      end

      # Fallback to filename if class not found
      class_name ||= File.basename(file_path, '.rb').capitalize

      file.puts "RSpec.describe #{class_name} do"
      file.puts "  # Your specs here"
      file.puts "end"
    end

    puts "Created spec file: #{spec_file_path}"
  end

  # Open the spec file
  system("zed #{spec_file_path}")
else
  puts "File must be in the app directory to create a spec for it"
  exit
end
