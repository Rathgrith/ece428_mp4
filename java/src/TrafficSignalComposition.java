import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import org.apache.hadoop.io.FloatWritable;
import org.apache.hadoop.mapreduce.lib.input.NLineInputFormat;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class TrafficSignalComposition {
    public static class TrafficSignalMapper extends Mapper<Object, Text, Text, Text> {
        private final Text interconneType = new Text();
        private final Text detectionType = new Text();

        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            String[] columns = value.toString().split(",");
            if (columns.length > 10) {
                String interconne = columns[10].trim();
                String detection = columns[9].trim();

                // Check if the 'Interconne' type matches the provided parameter 'X'
                if (interconne.equals(context.getConfiguration().get("interconneType"))) {
                    interconneType.set(interconne);
                    detectionType.set(detection);
                    context.write(interconneType, detectionType);
                }
            }
        }
    }

    public static class TrafficSignalReducer extends Reducer<Text, Text, Text, FloatWritable> {
        private final FloatWritable result = new FloatWritable();
        private float totalCount = 0;

        public void reduce(Text key, Iterable<Text> values, Context context)
                throws IOException, InterruptedException {
            Map<String, Integer> detectionCount = new HashMap<>();
            totalCount = 0;

            // Count the occurrences of each 'Detection_' value
            for (Text val : values) {
                String detection = val.toString();
                detectionCount.put(detection, detectionCount.getOrDefault(detection, 0) + 1);
                totalCount++;
            }

            // Calculate the percentage composition for each 'Detection_' value
            for (Map.Entry<String, Integer> entry : detectionCount.entrySet()) {
                String detection = entry.getKey();
                int count = entry.getValue();
                float composition = (count / totalCount) * 100;
                result.set(composition);
                context.write(new Text(detection), result);
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        conf.set("interconneType", args[2]); // Set the 'Interconne' type as a configuration parameter

        Job job = Job.getInstance(conf, "Traffic Signal Composition");
        job.setJarByClass(TrafficSignalComposition.class);
        job.setMapperClass(TrafficSignalMapper.class);
        job.setReducerClass(TrafficSignalReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class);
        job.setInputFormatClass(NLineInputFormat.class);
        NLineInputFormat.setNumLinesPerSplit(job, 14);
        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
