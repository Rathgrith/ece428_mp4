import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.FloatWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class TrafficSignalComposition {
    public static class TrafficSignalMapper extends Mapper<Object, Text, Text, FloatWritable> {
        private final Text detectionType = new Text();
        private final FloatWritable one = new FloatWritable(1);

        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            String[] columns = value.toString().split(",");
            if (columns.length > 10) {
                String interconne = columns[10].trim();
                String detection = columns[9].trim();

                // Check if the 'Interconne' type matches the provided parameter 'X'
                if (interconne.equals(context.getConfiguration().get("interconneType"))) {
                    detectionType.set(detection);
                    context.write(detectionType, one);
                }
            }
        }
    }

    public static class TrafficSignalReducer extends Reducer<Text, FloatWritable, Text, FloatWritable> {
        private final FloatWritable result = new FloatWritable();
        private float totalCount = 0;

        public void reduce(Text key, Iterable<FloatWritable> values, Context context)
                throws IOException, InterruptedException {
            float sum = 0;
            for (FloatWritable val : values) {
                sum += val.get();
                totalCount += val.get();
            }
            result.set((sum / totalCount) * 100); // Calculate the percentage composition
            context.write(key, result);
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        conf.set("interconneType", args[2]); // Set the 'Interconne' type as a configuration parameter

        Job job = Job.getInstance(conf, "Traffic Signal Composition");
        job.setJarByClass(TrafficSignalComposition.class);
        job.setMapperClass(TrafficSignalMapper.class);
        job.setCombinerClass(TrafficSignalReducer.class);
        job.setReducerClass(TrafficSignalReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(FloatWritable.class);

        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
